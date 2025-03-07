package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
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
	Plate     string
	Timestamp uint32
}

func ParseString(reader *bufio.Reader) (string, error) {
	// read the length byte
	lengthByte, err := reader.ReadByte()
	if err != nil {
		return "", fmt.Errorf("error reading plate length: %w", err)
	}
	stringLen := int(lengthByte)

	// read the string into the buffer
	stringBytes := make([]byte, stringLen)
	n, err := io.ReadFull(reader, stringBytes)
	if err != nil {
		return "", fmt.Errorf("error reading plate string: %w", err)
	}
	if n != stringLen {
		return "", fmt.Errorf("incomplete plate string: expected %d bytes, got %d", stringLen, n)
	}

	return string(stringBytes), nil
}

func ParseCameraRequest(reader *bufio.Reader) (Camera, error) {
	var data Camera
	err := binary.Read(reader, binary.BigEndian, &data)

	if err != nil {
		return data, fmt.Errorf("Error Parsing CameraRequest %x", err)
	}

	return data, nil
}

func ParsePlateRecord(reader *bufio.Reader) (Plate, error) {
	plateRecord := Plate{}
	// read the plate value
	plate, err := ParseString(reader)
	if err != err {
		return plateRecord, fmt.Errorf("error reading plate (str): %w", err)
	}
	plateRecord.Plate = plate

	// Read the timestamp (4 bytes).
	err = binary.Read(reader, binary.BigEndian, &plateRecord.Timestamp)
	if err != nil {
		return plateRecord, fmt.Errorf("error reading timestamp: %w", err)
	}

	return plateRecord, nil
}
