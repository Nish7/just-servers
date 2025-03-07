package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
)

type MsgType byte

const (
	CAMERA_REQ MsgType = 0x80
)

type CameraRequest struct {
	Road  uint16
	Mile  uint16
	Limit uint16
}

func HandleRequest[T any](reader *bufio.Reader, handler func(T) error) error {
	d, err := ParseRequest[T](reader)
	if err != nil {
		fmt.Println("Error Parsing CameraRequest", err)
		return nil
	}

	return handler(d)
}

func ParseRequest[T any](reader *bufio.Reader) (T, error) {
	var data T
	err := binary.Read(reader, binary.BigEndian, &data)
	if err != nil {
		return data, fmt.Errorf("failed to read request: %v", err)
	}
	return data, nil
}
