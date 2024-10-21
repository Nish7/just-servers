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
	buf := make(Request, 9)

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

		operation, n1, n2 := buf.Decode()
		fmt.Printf("Recieved (%s): %x (hex) - operation [%c], n2 [%d] n2 [%d]\n", conn.RemoteAddr().String(), buf, operation, n1, n2)
		s.HandleRequest(conn, operation, n1, n2)
	}
}

func (s *Server) HandleRequest(conn net.Conn, operation rune, n1 int32, n2 int32) {
	switch operation {
	case 'I':
		fmt.Print("Insert Operation!\n")
	case 'Q':
		fmt.Print("Query Operation!\n")
	default:
		fmt.Print("Invalid Operation\n")
	}
}
