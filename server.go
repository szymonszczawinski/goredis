package main

import (
	"fmt"
	"log/slog"
	"net"
)

const (
	defaultListenAddress = ":3333"
)

type Config struct {
	ListenAddress string
}

type Message struct {
	peer *Peer
	data []byte
}
type Server struct {
	peers       map[*Peer]bool
	addPeerChan chan *Peer
	quitChan    chan struct{}
	msgChan     chan Message
	listener    net.Listener
	kv          *ValKey
	Config
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddress) == 0 {
		cfg.ListenAddress = defaultListenAddress
	}
	return &Server{
		Config:      cfg,
		peers:       make(map[*Peer]bool),
		addPeerChan: make(chan *Peer),
		quitChan:    make(chan struct{}),
		msgChan:     make(chan Message),
		kv:          NewValKey(),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.ListenAddress)
	if err != nil {
		return err
	}
	s.listener = listener

	go s.loop()

	slog.Info("server runnin", "listen address", s.ListenAddress)
	return s.accept()
}

func (s *Server) accept() error {
	for {
		connection, err := s.listener.Accept()
		if err != nil {
			slog.Error("accept error", "error", err)
			continue
		}

		go s.handleConnection(connection)
	}
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgChan:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("raw message error", "err", err)
			}
		case <-s.quitChan:
			return
		case peer := <-s.addPeerChan:
			s.peers[peer] = true
		}
	}
}

func (s *Server) handleConnection(connection net.Conn) {
	peer := NewPeer(connection, s.msgChan)
	s.addPeerChan <- peer
	slog.Info("new peer connected", "remote address", connection.RemoteAddr().String())
	if err := peer.read(); err != nil {
		slog.Error("peer read error", "err", err, "remote address", connection.RemoteAddr().String())
	}
}

func (s *Server) handleMessage(message Message) error {
	cmd, err := ParseCommand(string(message.data))
	if err != nil {
		slog.Error("parse command error", "err", err)
		return err
	}
	slog.Info("parsed command", "cmd", cmd)
	switch v := cmd.(type) {
	case SetCommand:
		return s.kv.Set(string(v.key), v.value)
	case GetCommand:
		value, ok := s.kv.Get(string(v.key))
		if !ok {
			return fmt.Errorf("key not found %v", v.key)
		}
		_, err := message.peer.Send(value)
		if err != nil {
			slog.Error("peer send error", "err", err)
			return err
		}

	}
	return nil
}
