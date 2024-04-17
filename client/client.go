package client

import (
	"bytes"
	"context"
	"log"
	"log/slog"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	conn net.Conn
	addr string
}

func NewClinet(addr string) Client {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	return Client{
		addr: addr,
		conn: conn,
	}
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	slog.Info("clinet get", "key", key)
	var buff bytes.Buffer

	wr := resp.NewWriter(&buff)

	err := wr.WriteArray([]resp.Value{
		resp.StringValue("GET"),
		resp.StringValue(key),
	})
	if err != nil {
		slog.Error("error write array", "err", err)
	}
	_, err = c.conn.Write(buff.Bytes())
	if err != nil {
		return "", err
	}
	outBuf := make([]byte, 1024)

	n, err := c.conn.Read(outBuf)
	if err != nil {
		return "", err
	}
	return string(outBuf[:n]), nil
}

func (c *Client) Set(ctx context.Context, key, val string) error {
	slog.Info("clinet set", "key", key)
	var buff bytes.Buffer

	wr := resp.NewWriter(&buff)

	err := wr.WriteArray([]resp.Value{
		resp.StringValue("SET"),
		resp.StringValue(key),
		resp.StringValue(val),
	})
	if err != nil {
		slog.Error("error write array", "err", err)
	}
	_, err = c.conn.Write(buff.Bytes())
	return err
}
