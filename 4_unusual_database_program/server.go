package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Server struct {
	quitch chan struct{}
	addr   string
	store  map[string]string
}

func NewServer(addr string) *Server {
	return &Server{
		quitch: make(chan struct{}),
		addr:   addr,
		store: map[string]string{
			"version": "1.0",
		},
	}
}

func (s *Server) Start() error {
	addr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		Log("error: resolving from udp addr"+err.Error(), addr)
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		Log("error: listening on the port "+err.Error(), addr)
		return err
	}

	defer conn.Close()

	buf := make([]byte, 1000)

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			if err == io.EOF {
				Log("EOF Detected: "+err.Error(), addr)
				break
			}

			Log("Error:in reading from udp "+err.Error(), addr)
			break
		}

		req := strings.Trim(string(buf[:n]), "\n")
		Log(fmt.Sprintf("Request: %s\n", req), addr)

		if before, after, found := strings.Cut(req, "="); found {
			Log(fmt.Sprintf("Insert Operation: key : %s and value %s", before, after), addr)
			err := s.InsertOp(before, after)

			if err != nil {
				Log("InsertionError: "+err.Error(), addr)
				continue
			}

			Log("Inserted: "+before, addr)
		} else {
			Log(fmt.Sprintf("Retrieve Operation: %s", before), addr)
			val, err := s.RetrieveOp(req)

			if err != nil {
				res := req + "=" + "\n"
				Log(fmt.Sprintf("Error: %e - Response: %s", err, res), addr)
				conn.WriteTo([]byte(res), addr)
				continue
			}

			res := req + "=" + val + "\n"
			Log("Response: "+res, addr)
			conn.WriteTo([]byte(res), addr)
		}
	}

	return nil
}

func (s *Server) InsertOp(key, val string) error {
	if key == "version" {
		return errors.New("You cannot update version!")
	}

	s.store[key] = val
	return nil
}

func (s *Server) RetrieveOp(key string) (string, error) {
	val, ok := s.store[key]

	if !ok {
		return "", errors.New("Missing Key: " + key)
	}

	return val, nil
}

func Log(p string, addr net.Addr) {
	log.Printf("[%s] %s", addr.String(), p)
}
