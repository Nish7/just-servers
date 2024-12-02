package main

import "log"

func main() {
	server := NewServer(":8082")
	err := server.StartServer()
	if err != nil {
		log.Fatalf("Error: Starting the server %v", err)
	}
}
