package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
)

type Request struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func handleRequest(c net.Conn) {
	// handle decoding
	defer c.Close()
	reader := bufio.NewReader(c)

	for {
		message, err := reader.ReadString('\n')

		if err != nil {
			c.Write([]byte(err.Error()))
			break
		}

		fmt.Printf("Recieved: %v", message)

		var request Request
		err = json.Unmarshal([]byte(message), &request)

		// handle validation
		if err != nil {
			c.Write([]byte(err.Error()))
			break
		}

		if request.Number < 1 {
			c.Write([]byte("Number should be greater 1"))
			break
		}

		if request.Method != "isPrime" {
			c.Write([]byte("Request Method is not isPrime"))
			break
		}

		// handle response
		res := Response{Method: "isPrime", Prime: isPrime(request.Number)}
		resJson, err := json.Marshal(res)

		if err != nil {
			log.Fatal(err)
		}

		c.Write([]byte(string(resJson) + "\n"))
	}
}

func isPrime(n int) bool {
	if n == 1 {
		return false
	}

	if n == 2 {
		return true
	}

	if n%2 == 0 {
		return false
	}

	sqrtN := int(math.Sqrt(float64(n)))
	for i := 3; i <= sqrtN; i += 2 {
		if n%i == 0 {
			return false
		}
	}

	return true
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
		conn, err := l.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go handleRequest(conn)
	}

}
