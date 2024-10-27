package main

import (
	"log"
	"net"
)

type Server struct {
	quitch   chan struct{}
	listener net.Listener
	addr     string
	users    map[string]net.Conn
}

func NewServer(addr string) *Server {
	return &Server{
		quitch: make(chan struct{}),
		addr:   addr,
		users:  make(map[string]net.Conn),
	}
}

func (s *Server) Start(addr string) error {
	l, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}

	log.Printf("Server Listening on Port %s", addr)
	s.listener = l
	go s.Accept()

	<-s.quitch
	defer l.Close()
	return nil
}

func (s *Server) Accept() {
	for {
		conn, err := s.listener.Accept()
		log.Printf("New Connection: %s\n", conn.RemoteAddr().String())

		if err != nil {
			log.Printf("connection error: %v\n", err)
			continue
		}

		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	// handle conn
}
