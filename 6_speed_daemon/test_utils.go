package main

import (
	"bufio"
	"testing"
	"time"
)

func SetupDispatchers(t *testing.T, addr string, dispatchers ...Dispatcher) []*TCPClient {
	t.Helper()

	clients := make([]*TCPClient, len(dispatchers))
	for i, dispatcher := range dispatchers {
		client := NewTCPClient(addr)
		if err := client.Connect(); err != nil {
			t.Fatalf("Failed to connect dispatcher %d: %v", i, err)
		}

		client.SendIAMDispatcher(dispatcher)
		clients[i] = client
	}
	return clients
}

func SetupCameras(t *testing.T, addr string, cameras ...Camera) []*TCPClient {
	t.Helper()

	clients := make([]*TCPClient, len(cameras))
	for i, camera := range cameras {
		client := NewTCPClient(addr)
		if err := client.Connect(); err != nil {
			t.Fatalf("Failed to connect camera %d: %v", i, err)
		}

		client.SendIAMCamera(camera)
		clients[i] = client
	}
	return clients
}

func ClientCleanUp(t *testing.T, clients ...*TCPClient) {
	t.Helper()
	for _, client := range clients {
		if client != nil {
			client.Disconnect()
		}
	}
}

func AssertError(t *testing.T, reader *bufio.Reader, expectedError ErrorResp) {
	t.Helper()

	msgType, _ := reader.ReadByte()
	if msgType != byte(ERROR_RESP) {
		t.Fatalf("Illegal Message Type/Code")
	}

	recievedError, err := ParseError(reader)
	if err != nil {
		t.Fatalf("Error Parsing Ticket. Wrong Message")
	}

	if recievedError != expectedError {
		t.Fatalf("ExpectedError %v and Got %v", expectedError, recievedError)
	}

}

func AssertTicket(t *testing.T, reader *bufio.Reader, expectedTicket Ticket) {
	t.Helper()

	msgType, _ := reader.ReadByte()
	if msgType != byte(TICKET_RESP) {
		t.Fatalf("Illegal Message Type/Code")
	}

	recievedTicket, err := ParseTicket(reader)
	if err != nil {
		t.Fatalf("Error Parsing Ticket. Wrong Message")
	}

	if recievedTicket != expectedTicket {
		t.Fatalf("ExpectedTicket %v and Got %v", expectedTicket, recievedTicket)
	}
}

func AssertHeartbeat(t *testing.T, reader *bufio.Reader, expectedInterval uint32) {
	t.Helper()

	start := time.Now()
	msgType, _ := reader.ReadByte()
	if msgType != byte(HEARTBEAT_RESP) {
		t.Fatalf("Illegal Message Type/Code: got %v, want %v", msgType, HEARTBEAT_RESP)
	}

	elapsed := time.Since(start)
	expected := time.Duration(expectedInterval*100) * time.Millisecond
	if elapsed < expected-50*time.Millisecond || elapsed > expected+50*time.Millisecond {
		t.Fatalf("Heartbeat interval mismatch: got %v, want ~%v", elapsed, expected)
	}
}

func AssertNoHeartbeat(t *testing.T, reader *bufio.Reader, timeout time.Duration) {
	t.Helper()

	done := make(chan struct{})
	var receivedByte byte

	go func() {
		b, err := reader.ReadByte()
		if err == nil {
			receivedByte = b
		}
		close(done)
	}()

	select {
	case <-time.After(timeout):
		return
	case <-done:
		if receivedByte == byte(HEARTBEAT_RESP) {
			t.Fatalf("Received unexpected heartbeat message: got type %v", receivedByte)
		} else {
			t.Fatalf("Received unexpected message: got type %v, expected no messages", receivedByte)
		}
	}
}
