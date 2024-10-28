package main

import "log"

func main() {
	server := NewServer(":8080")
	err := server.Start()

	if err != nil {
		log.Fatalf("Error during listenting [%v]", err)
	}
}
