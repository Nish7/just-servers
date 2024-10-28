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

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	s.presenceNotification(conn)
	s.userMap.AddUser(nickname, conn)

	// TODO: User Join Broadcast

	// TODO: defer leave user
	log.Printf("%s joined the room", nickname)

	for scanner.Scan() {
		log.Printf("Recieved: %s", scanner.Text())
	}
}

func (s *Server) joinRequest(scanner *bufio.Scanner) (nickname string, err error) {
	ok := scanner.Scan()

	if !ok {
		return "", scanner.Err()
	}

	inputName := scanner.Text()

	// check the length of the name
	if len(inputName) < 1 && len(inputName) > 18 {
		return "", errors.New("Length of the name is less than 1 or greater than 18")
	}

	// check if the name contains only letters and numbers
	for _, r := range inputName {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return "", errors.New("Invalid Characters")
		}
	}

	// check if the name is already taken
	if _, ok := s.userMap.getConnection(nickname); ok {
		return "", errors.New("Name already taken")
	}

	return inputName, nil
}

func (s *Server) presenceNotification(conn net.Conn) {
	roomMembers := s.userMap.GetNicknames()

	if len(roomMembers) > 0 {
		nicknames := strings.Join(roomMembers, ", ")
		conn.Write([]byte("> the room contains: " + nicknames + "\n"))
	}
}
