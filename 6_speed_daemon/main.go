package main

import "log"

func main() {
	server := NewServer(":8085", NewInMemoryStore())
	err := server.Start()
	if err != nil {
		log.Fatalf("Error: Starting the server %v", err)
	}
}
