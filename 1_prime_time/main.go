package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func isPrime(n int) bool {
	return true
}

type RequestResponse struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}

func handleRequest(c net.Conn) {
	defer c.Close()
	// handle decoding
	reader := bufio.NewReader(c)
	message, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	log.Println(message)

	var request RequestResponse
	err = json.Unmarshal([]byte(message), &request)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Recieved: %v\n", request)

	// handle malformed validaton
}

func main() {
	// create a listener
	l, err := net.Listen("tcp", ":8081")

	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	fmt.Println("Server is listening on port 8081...")

	for {
		// accept a connection
		conn, err := l.Accept()

		if err != nil {
			log.Fatal(err)
		}

		// use a go routine to handle the response
		go handleRequest(conn)
	}

}
