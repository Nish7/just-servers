package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"unicode"
)

const (
	CHAT_SERVER_URL  = "chat.protohackers.com"
	CHAT_SERVER_PORT = "16963"
	BOGUS_COIN_ADDR  = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
)

type Server struct {
	quitch   chan struct{}
	listener net.Listener
	addr     string
}

func NewServer(addr string) *Server {
	return &Server{
		addr:   addr,
		quitch: make(chan struct{}),
	}
}

func (s *Server) StartServer() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.listener = l
	fmt.Printf("Server Listening on the Port %s\n", s.addr)

	go s.Accept()

	<-s.quitch
	defer l.Close()
	return nil
}

func (s *Server) Accept() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			Log(fmt.Sprintf("Connection Error: %v\n", err), conn.RemoteAddr())
		}

		Log("Client Connected!\n", conn.RemoteAddr())
		go s.HandleClientConn(conn)
	}
}

func (s *Server) HandleClientConn(conn net.Conn) {
	defer conn.Close()

	// connect to the upstream server
	serverconn, err := net.Dial("tcp", CHAT_SERVER_URL+":"+CHAT_SERVER_PORT)
	if err != nil {
		Log(fmt.Sprintf("Connection Error with Chat Server: %v\n", err), conn.RemoteAddr())
		return
	}

	Log("Connected to the Chat Server\n", conn.RemoteAddr())
	go s.HandleServerConn(serverconn, conn)
	defer serverconn.Close()

	// send all incoming messages from client to the chat server
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())
			} else {
				fmt.Printf("Error reading from client %s: %s\n", conn.RemoteAddr(), err)
			}
			break
		}

		Log(fmt.Sprintf("Recieved from Client: %q\n", line), conn.RemoteAddr())
		text := rewriteAddr(line)

		_, err = serverconn.Write([]byte(text))
		if err != nil {
			Log(fmt.Sprintf("Write Error: cannot send to the server %v\n", err), conn.RemoteAddr())
		}

		Log(fmt.Sprintf("Sent: %q\n", text), conn.RemoteAddr())
	}
}

func (s *Server) HandleServerConn(serverconn net.Conn, clientconn net.Conn) {
	// send all incoming messages from server to the client
	scanner := bufio.NewScanner(serverconn)

	for scanner.Scan() {
		Log(fmt.Sprintf("Recieved from Server: %q\n", scanner.Text()), clientconn.RemoteAddr())
		text := rewriteAddr(scanner.Text())

		_, err := clientconn.Write([]byte(text + "\n"))
		if err != nil {
			Log(fmt.Sprintf("Write Error: cannot send to the client %v", err), clientconn.RemoteAddr())
		}
	}
}

// Malicious Function
func rewriteAddr(message string) string {
	words := strings.Split(message, " ")

	for _, word := range words {
		trimmedWord := strings.TrimSuffix(word, "\n")
		if isBogusCoinAddr(trimmedWord) {
			message = strings.ReplaceAll(message, trimmedWord, BOGUS_COIN_ADDR)
		}
	}

	return message
}

func isBogusCoinAddr(word string) bool {
	word = strings.TrimSuffix(word, "\n")
	if len(word) < 27 || len(word) > 35 || word[0] != '7' {
		return false
	}

	// check if alphanumeric
	for _, ch := range word {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
			return false
		}
	}

	return true
}

func Log(p string, addr net.Addr) {
	log.Printf("[%s] %s", addr.String(), p)
}
