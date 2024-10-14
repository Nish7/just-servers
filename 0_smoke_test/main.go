package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	fmt.Printf("Listening on port 8080")
	l, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	for {
		// wait for the conn
		conn, err := l.Accept() // returns a single connection to a remote addr

		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) { // run a evey request in a go routine.
			io.Copy(c, c) // connection implemented reader and writer, thus basically reading and writing to the same buffer
			c.Close()
		}(conn)
	}
}
