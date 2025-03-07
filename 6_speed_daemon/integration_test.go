package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"
)

var testServer *Server
var addr string = ":8080"

func TestMain(t *testing.M) {
	// setup the server
	testServer = NewServer(addr)

	go func() {
		err := testServer.Start()
		if err != nil {
			log.Fatalf("Error: Starting the server %v", err)
		}
	}()

	time.Sleep(1000 * time.Millisecond) // give some time to start the server
	code := t.Run()
	os.Exit(code)
}

func TestCameraRequest(t *testing.T) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Test multiple IAmCamera messages
	cameras := []struct{ road, mile, limit uint16 }{
		{66, 100, 60},
		{123, 8, 60},
		{368, 1234, 40},
	}

	for _, cam := range cameras {
		msg := make([]byte, 7)
		msg[0] = 0x80
		binary.BigEndian.PutUint16(msg[1:3], cam.road)
		binary.BigEndian.PutUint16(msg[3:5], cam.mile)
		binary.BigEndian.PutUint16(msg[5:7], cam.limit)

		_, err = conn.Write(msg)

		if err != nil {
			fmt.Println("Error sending IAmCamera message")
			return
		}

		fmt.Printf("Client -> Sent CameraRequest: %x\n", msg)
	}

	time.Sleep(500 * time.Millisecond) // test ended before verifying
}
