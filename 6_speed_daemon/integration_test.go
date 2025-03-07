package main

import (
	"log"
	"os"
	"testing"
	"time"
)

var testServer *Server
var addr string = ":8080"

func TestMain(t *testing.M) {
	// setup the server
	store := NewInMemoryStore()
	testServer = NewServer(addr, store)

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

func TestPlateRequest(t *testing.T) {
	client := NewTCPClient(addr)
	client.Connect()
	defer client.Disconnect()

	client.SendIAMCamera(Camera{20, 80, 100})
	client.SendPlateRecord(Plate{"UN1X", 1000})

	time.Sleep(500 * time.Millisecond) // test ended before verifying
}

func TestDispatcherRequest(t *testing.T) {
	client := NewTCPClient(addr)
	client.Connect()
	defer client.Disconnect()

	client.SendIAMDispatcher(Dispatcher{[]uint16{66}})
	time.Sleep(500 * time.Millisecond) // test ended before verifying
}

func TestCameraRequest(t *testing.T) {
	client := NewTCPClient(addr)
	client.Connect()
	defer client.Disconnect()

	client.SendIAMCamera(Camera{66, 100, 60})
	time.Sleep(500 * time.Millisecond) // test ended before verifying
}
