package main

import (
	"bufio"
	"errors"
	"log"
	"net"
	"strings"
	"unicode"
)

type Server struct {
	quitch   chan struct{}
	listener net.Listener
	addr     string
	userMap  *UserMap
}

func NewServer(addr string) *Server {
	return &Server{
		quitch:  make(chan struct{}),
		addr:    addr,
		userMap: NewUsersMap(),
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
		log.Printf("New Connection: %s\n", conn.RemoteAddr().String())

		if err != nil {
			log.Printf("connection error: %v\n", err)
			continue
		}

		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))
	scanner := bufio.NewScanner(conn)
	nickname, err := s.joinRequest(scanner)

	if err != nil || nickname == "" {
		log.Printf("validation error: %v", err)
		return
	}

	s.presenceNotification(conn)
	s.userMap.AddUser(nickname, conn)

	defer s.Leave(nickname)

	s.Broadcast(nickname, "* "+nickname+" joined the room")
	log.Printf("%s joined the room", nickname)

	for scanner.Scan() {
		message := strings.TrimSpace(scanner.Text())

		if len(message) > 1000 {
			conn.Write([]byte("message is too long. Re-send the message\n"))
			continue
		}

		s.Broadcast(nickname, "["+nickname+"] "+message)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("scanner error: %v", err)
	}
}

func (s *Server) joinRequest(scanner *bufio.Scanner) (string, error) {
	ok := scanner.Scan()

	if !ok {
		return "", scanner.Err()
	}

	inputName := strings.TrimSpace(scanner.Text())
	log.Print("Received name: ", inputName)

	if len(inputName) < 1 || len(inputName) > 18 {
		return "", errors.New("length of the name is less than 2 or greater than 19")
	}

	for _, r := range inputName {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return "", errors.New("invalid characters")
		}
	}

	if _, ok := s.userMap.getConnection(inputName); ok {
		return "", errors.New("name already taken")
	}

	return inputName, nil
}

func (s *Server) presenceNotification(conn net.Conn) {
	roomMembers := s.userMap.GetNicknames()
	nicknames := strings.Join(roomMembers, ", ")
	conn.Write([]byte("* the room contains: " + nicknames + "\n"))
}

func (s *Server) Broadcast(sender string, message string) {
	for _, key := range s.userMap.GetNicknames() {
		if key != sender {
			conn, _ := s.userMap.getConnection(key)
			conn.Write([]byte(message + "\n"))
		}
	}
}

func (s *Server) Leave(nickname string) {
	s.userMap.RemoveUser(nickname)
	s.Broadcast(nickname, "* "+nickname+" left the room")
	log.Printf("%s left the room", nickname)
}
