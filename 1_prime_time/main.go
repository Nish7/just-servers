package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
)

type Request struct {
	Method *string    `json:"method"`
	Number *BigNumber `json:"number"`
}

func (req *Request) validFields() bool {
	if req.Method == nil {
		return false
	}

	if *req.Method != "isPrime" {
		return false
	}

	if req.Number == nil {
		return false
	}

	return true
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
			log.Print(err.Error())
			c.Write([]byte(err.Error()))
			break
		}

		log.Printf("Recieved: %v", message)

		var request Request
		err = json.Unmarshal([]byte(message), &request)

		// handle validation
		if err != nil {
			log.Print(err.Error())
			c.Write([]byte(err.Error()))
			break
		}

		if !request.validFields() {
			log.Print("Invalid Fields")
			c.Write([]byte("Invalid Fields"))
			break
		}

		// convert the string to bigInt
		if request.Number.IsFloat {
			res := Response{Method: "isPrime", Prime: false}
			handleResponse(c, res)
			break
		}

		// handle response
		res := Response{Method: "isPrime", Prime: isPrime(*request.Number.BigInt)}
		handleResponse(c, res)
	}
}

func handleResponse(c net.Conn, res Response) {
	resJson, err := json.Marshal(res)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Response: ", string(resJson)+"\n")
	c.Write([]byte(string(resJson) + "\n"))
}

func isPrime(n big.Int) bool {
	k := 10 // A higher k increases the confidence that the number is prime, but it also takes more time.
	return n.ProbablyPrime(k)
}

func main() {
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
