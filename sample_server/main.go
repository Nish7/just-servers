package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{} // signal to gracefully shut down the server from any go routine
	msgch      chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	fmt.Printf("Server Listening on %s\n", s.listenAddr)
	if err != nil {
		return err
	}

	defer ln.Close()
	s.ln = ln

	go s.loop()

	<-s.quitch
	close(s.msgch)
	return nil
}

func (s *Server) loop() {
	for { // keep acepting new connections and perform connection level processing in a new go routine
		conn, err := s.ln.Accept()
		fmt.Printf("new connection: %s", conn.RemoteAddr())

		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}

		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)

	for { // keep taking request from the same connection
		n, err := conn.Read(buf)

		if err != nil {
			fmt.Println("read error:", err)
			break
		}

		s.msgch <- Message{from: conn.RemoteAddr().String(), payload: buf[:n]}

		conn.Write([]byte("Thank you for the message!\n"))
	}
}

func (s *Server) HandleMessage() {
	for msg := range s.msgch {
		fmt.Printf("recieved message from the connection (%s): %s\n", msg.from, string(msg.payload))
	}
}

func main() {
	server := NewServer(":3001")
	go server.HandleMessage()
	log.Fatal(server.Start())
}
