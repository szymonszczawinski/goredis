package main

import (
	"log/slog"
	"net"
)

type Peer struct {
	connection net.Conn
}

func NewPeer(connection net.Conn) *Peer {
	return &Peer{
		connection: connection,
	}
}

func (p *Peer) read() error {
	buffer := make([]byte, 1024)
	for {
		bytesRead, err := p.connection.Read(buffer)
		if err != nil {
			slog.Error("peer read error", "err", err)
			return err
		}
		msg := buffer[:bytesRead]
	}
}
