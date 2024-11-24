package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Listening on port 8082")
	s := NewServer(":8082")

	if s.Start() != nil {
		log.Fatalf("Error: creating a server")
	}
}
