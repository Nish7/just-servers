package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

type Server struct {
	quitch   chan struct{}
	listener net.Listener
	addr     string
}

func NewServer(addr string) *Server {
	return &Server{
		quitch: make(chan struct{}),
		addr:   addr,
	}
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.addr)

	if err != nil {
		return err
	}

	log.Printf("Server Listening on Port %s", s.addr)
	s.listener = l
	go s.Accept()

	<-s.quitch
	defer l.Close()
	return nil
}

func (s *Server) Accept() {
	for {
		conn, err := s.listener.Accept()
		log.Printf("New Connection :%s\n", conn.RemoteAddr().String())

		if err != nil {
			log.Printf("Error: Connection error %e\n", err)
			return
		}

		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		msgType, err := reader.ReadByte()
		if err == io.EOF {
			log.Println("Connection closed by remote end")
			return
		}

		if err != nil {
			log.Printf("Error Reading Connection %e", err)
			return
		}

		// request handler
		switch msgType {
		case 0x80:
			err = s.HandleCameraRequest(reader)
			if err != nil {
				fmt.Println("Error handling CameraRequest", err)
				return
			}
		default:
			fmt.Printf("Unknown message type: %x echoing back\n", msgType)
		}
	}
}

func (s *Server) HandleCameraRequest(reader *bufio.Reader) error {
	// Expect 6 more byte: 3 x u16 (road, mile, limit)
	data := make([]byte, 6)
	n, err := reader.Read(data)

	if err != nil {
		return fmt.Errorf("failed to read Camera Request %v", err)
	}

	if n != 6 {
		return fmt.Errorf("imcomplete CameraRequest %v", err)
	}

	// parse the big-endian u16 fields
	road := binary.BigEndian.Uint16(data[0:2])
	mile := binary.BigEndian.Uint16(data[2:4])
	limit := binary.BigEndian.Uint16(data[4:6])

	fmt.Printf("CameraRequest Received: Road=%d, mile=%d, limit=%d\n", road, mile, limit)
	return nil
}
