package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type Server struct {
	quitch        chan struct{}
	listener      net.Listener
	addr          string
	store         Store
	cameras       map[net.Conn]Camera
	dispatchers   map[net.Conn]Dispatcher
	pending_queue []Ticket
	slock         sync.Mutex
}

func NewServer(addr string, store Store) *Server {
	return &Server{
		quitch:        make(chan struct{}),
		addr:          addr,
		store:         store,
		cameras:       make(map[net.Conn]Camera),
		dispatchers:   make(map[net.Conn]Dispatcher),
		pending_queue: []Ticket{},
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
	reader := bufio.NewReader(conn)
	clientType := UNKNOWN
	heartbeatRegistered := false

	defer s.CleanUpClient(conn, &clientType)

	for {
		msgType, err := ReadMsgType(reader)
		if err != nil {
			log.Printf("[%s] Connection Error: %v", conn.RemoteAddr().String(), err)
			return
		}

		switch msgType {
		case IAMCAMERA_REQ:
			err = s.HandleCameraReq(conn, reader, &clientType)
		case IAMDISPATCHER_REQ:
			err = s.HandleDispatcherReq(conn, reader, &clientType)
		case PLATE_REQ:
			err = s.HandlePlateReq(conn, reader, &clientType)
		case WANTHEARTBEAT_REQ:
			err = s.WantHeatbeatHandler(conn, reader, &heartbeatRegistered, &clientType)
		default:
			err = fmt.Errorf("Unknown message type: %X\n", msgType)
		}

		if err != nil {
			s.ErrorHandler(err, conn)
			return
		}
	}
}

func (s *Server) CleanUpClient(conn net.Conn, client *ClientType) {
	defer conn.Close()
	switch *client {
	case CAMERA:
		s.slock.Lock()
		delete(s.cameras, conn)
		s.slock.Unlock()
	case DISPATCHER:
		s.slock.Lock()
		delete(s.dispatchers, conn)
		s.slock.Unlock()
	}
}

func ReadMsgType(reader *bufio.Reader) (MsgType, error) {
	msgType, err := reader.ReadByte()
	if err == io.EOF {
		return 0, fmt.Errorf("Connection closed by remote end")
	}

	if err != nil {
		return 0, fmt.Errorf("Unknown Error %v", err)
	}

	return MsgType(msgType), nil
}
