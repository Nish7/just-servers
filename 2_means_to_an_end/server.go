package main

import (
	"fmt"
	"io"
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
	fmt.Printf("Server Listening on %s\n", s.addr)
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
		fmt.Printf("New Connection: %s\n", conn.RemoteAddr().String())

		if err != nil {
			fmt.Printf("connection error: %v\n", err)
			continue
		}

		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 9)

	for {
		_, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				fmt.Printf("Connection closed by client: %s\n", conn.RemoteAddr().String())
			} else {
				fmt.Printf("Read error: %v\n", err)
			}
			break
		}

		fmt.Printf("Recieved (%s): %x\n", conn.RemoteAddr().String(), buf)
	}
}
