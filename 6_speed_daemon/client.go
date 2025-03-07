package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type TCPClient struct {
	address string
	conn    net.Conn
}

func NewTCPClient(address string) *TCPClient {
	return &TCPClient{
		address: address,
	}
}

func (c *TCPClient) Connect() error {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", c.address, err)
	}
	c.conn = conn
	return nil
}

func (c *TCPClient) SendPlateRecord(plate Plate) {
	plateBytes := []byte(plate.plate)
	plateLen := len(plateBytes)

	msg := make([]byte, 1+1+plateLen+4)

	msg[0] = byte(PLATE_REQ)
	msg[1] = byte(plateLen)

	copy(msg[2:2+plateLen], plateBytes)
	binary.BigEndian.PutUint32(msg[2+plateLen:], plate.timestamp)

	_, err := c.conn.Write(msg)
	if err != nil {
		fmt.Printf("Error sending PlateRecord message: %v\n", err)
		return
	}

	fmt.Printf("Client -> Sent PlateRecord: %v - hex[% X]\n", plate, msg)
}

func (c *TCPClient) SendIAMCamera(cam Camera) {
	msg := make([]byte, 7)
	msg[0] = byte(IAMCAMERA_REQ)
	binary.BigEndian.PutUint16(msg[1:3], cam.Road)
	binary.BigEndian.PutUint16(msg[3:5], cam.Mile)
	binary.BigEndian.PutUint16(msg[5:7], cam.Limit)

	_, err := c.conn.Write(msg)

	if err != nil {
		fmt.Printf("Error sending IAmCamera message")
		return
	}

	fmt.Printf("Client -> Sent CameraRequest: %v - hex[% X]\n", cam, msg)
}

func (c *TCPClient) Disconnect() {
	c.conn.Close()
}
