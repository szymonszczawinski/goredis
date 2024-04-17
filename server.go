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
	cmd  Command
}
type Server struct {
	peers          map[*Peer]bool
	addPeerChannel chan *Peer
	quitChannel    chan struct{}
	messageChannel chan Message
	listener       net.Listener
	storage        *ValKey
	Config
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddress) == 0 {
		cfg.ListenAddress = defaultListenAddress
	}
	return &Server{
		Config:         cfg,
		peers:          make(map[*Peer]bool),
		addPeerChannel: make(chan *Peer),
		quitChannel:    make(chan struct{}),
		messageChannel: make(chan Message),
		storage:        NewValKey(),
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
		case message := <-s.messageChannel:
			if err := s.handleMessage(message); err != nil {
				slog.Error("raw message error", "err", err)
			}
		case <-s.quitChannel:
			return
		case peer := <-s.addPeerChannel:
			s.peers[peer] = true
		}
	}
}

func (s *Server) handleConnection(connection net.Conn) {
	peer := NewPeer(connection, s.messageChannel)
	s.addPeerChannel <- peer
	slog.Info("new peer connected", "remote address", connection.RemoteAddr().String())
	if err := peer.ReadLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remote address", connection.RemoteAddr().String())
	}
}

func (s *Server) handleMessage(message Message) error {
	switch v := message.cmd.(type) {
	case SetCommand:
		return s.storage.Set(string(v.key), v.value)
	case GetCommand:
		value, ok := s.storage.Get(string(v.key))
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
