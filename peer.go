package main

import (
	"log/slog"
	"net"
)

type Peer struct {
	connection net.Conn
	msgChan    chan Message
}

func NewPeer(connection net.Conn, msgChan chan Message) *Peer {
	return &Peer{
		connection: connection,
		msgChan:    msgChan,
	}
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.connection.Write(msg)
}

func (p *Peer) read() error {
	buffer := make([]byte, 1024)
	for {
		bytesRead, err := p.connection.Read(buffer)
		if err != nil {
			slog.Error("peer read error", "err", err)
			return err
		}
		msgBuf := make([]byte, bytesRead)
		copy(msgBuf, buffer)
		p.msgChan <- Message{peer: p, data: msgBuf}
	}
}
