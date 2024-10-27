package main

import "log"

func main() {
	server := NewServer("8084")
	err := server.Start()

	if err != nil {
		log.Fatalf("Error during listenting [%v]", err)
	}
}
