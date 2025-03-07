package main

import (
	"bufio"
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
		switch MsgType(msgType) {
		case CAMERA_REQ:
			err := HandleRequest(reader, s.handleCameraRequest)
			if err != nil {
				fmt.Printf("Error Handling Request: %x", err)
			}
		default:
			fmt.Printf("Unknown message type: %x\n", msgType)
		}
	}
}

func (s *Server) handleCameraRequest(req CameraRequest) error {
	fmt.Printf("Handling Camera Request %v\n", req)
	return nil
}
