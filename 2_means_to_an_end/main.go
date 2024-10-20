package main

func main() {
	server := NewServer(":8082")
	server.Start()
}
