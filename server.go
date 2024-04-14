package main

import (
	"log/slog"
	"net"
)

const (
	defaultListenAddress = ":3333"
)

type Config struct {
	ListenAddress string
}
type Server struct {
	peers       map[*Peer]bool
	addPeerChan chan *Peer
	quitChan    chan struct{}
	msgChan     chan []byte
	listener    net.Listener
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
		msgChan:     make(chan []byte),
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
		case rawMsg := <-s.msgChan:
			if err := s.handleRawMessage(rawMsg); err != nil {
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

func (s *Server) handleRawMessage(rawMessage []byte) error {
	slog.Info("handle raw message", "msg", string(rawMessage))
	return nil
}
