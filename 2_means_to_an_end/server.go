package main

import (
	"fmt"
	"net"
)

type Server struct {
	quitch   chan struct{}
	listener net.Listener
	addr     string
}

func NewServer(addr string) *Server {
	return &Server{
		quitch: make(chan struct{}, 1024),
		addr:   addr,
	}
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.addr)
	s.listener = l

	if err != nil {
		return err
	}

	go s.Accept(l)

	<-s.quitch
	defer l.Close()
	return nil
}

func (s *Server) Accept(l net.Listener) {
	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Printf("connection error: %v\n", err)
			continue
		}

		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		// handle the request
	}
}
