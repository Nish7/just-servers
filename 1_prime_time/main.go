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
	Number *big.Float `json:"number"`
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

		// handle response
		res := Response{Method: "isPrime", Prime: isPrime(*request.Number)}
		resJson, err := json.Marshal(res)

		if err != nil {
			log.Fatal(err)
		}

		log.Print("Response: ", string(resJson)+"\n")
		c.Write([]byte(string(resJson) + "\n"))
	}
}

func isPrime(n big.Float) bool {
	// convert to a big int. if a float than return false
	if !n.IsInt() {
		return false
	}

	bigInt := new(big.Int)
	n.Int(bigInt)
	k := 10 // A higher k increases the confidence that the number is prime, but it also takes more time.

	return bigInt.ProbablyPrime(k)
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
