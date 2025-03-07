package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
)

type MsgType byte

const (
	IAMCAMERA_REQ MsgType = 0x80
	PLATE_REQ     MsgType = 0x20
)

type Client int

const (
	CAMERA Client = iota
)

type Camera struct {
	Road  uint16
	Mile  uint16
	Limit uint16
}

type Plate struct {
	len       uint8
	plate     string
	timestamp uint32
}

func ParseRequest[T any](reader *bufio.Reader) (T, error) {
	var data T
	err := binary.Read(reader, binary.BigEndian, &data)

	if err != nil {
		return data, fmt.Errorf("Error Parsing CameraRequest %x", err)
	}

	return data, nil
}
