package main

import (
	"encoding/binary"
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
	request := make(Request, 9)
	store := NewInMemoryStore()

	for {
		err := binary.Read(conn, binary.BigEndian, request)
		fmt.Printf("Binary Data %b\n", request)

		if err != nil {
			if err == io.EOF {
				fmt.Printf("Connection closed by client: %s\n", conn.RemoteAddr().String())
			} else {
				fmt.Printf("Read error: %v\n", err)
			}
			break
		}

		operation, n1, n2 := request.Decode()
		fmt.Printf("Recieved (%s): %x (hex) - operation [%c], n2 [%d] n2 [%d]\n", conn.RemoteAddr().String(), request, operation, n1, n2)
		s.HandleRequest(store, conn, operation, n1, n2)
	}
}

func (s *Server) HandleRequest(store Store, conn net.Conn, operation rune, n1 int32, n2 int32) {
	switch operation {
	case 'I':
		s.HandleInsert(store, n1, n2)
	case 'Q':
		s.HandleQuery(store, conn, n1, n2)
	default:
		conn.Write([]byte("Invalid Operation"))
		fmt.Print("Invalid Operation\n")
	}
}

func (s *Server) HandleInsert(store Store, timestamp, price int32) {
	store.Insert(timestamp, price)
	fmt.Printf("Value Inserted:\n Store = %v\n", store)
}

func (s *Server) HandleQuery(store Store, conn net.Conn, minTime, maxTime int32) {
	mean := store.Query(minTime, maxTime)

	response := make([]byte, 4)
	binary.BigEndian.PutUint32(response, uint32(mean))
	fmt.Printf("For the meanTime [%d] and maxTime [%d]. The mean price is = [%d] (dec) - [%x] (hex)\n", minTime, maxTime, mean, response)
	conn.Write(response)
}
