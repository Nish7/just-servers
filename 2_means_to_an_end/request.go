package main

import "encoding/binary"

type Request []byte

func (r *Request) Decode() (operation rune, n1 int32, n2 int32) {
	operation = rune((*r)[0])
	n1 = int32(binary.BigEndian.Uint32((*r)[1:5]))
	n2 = int32(binary.BigEndian.Uint32((*r)[5:9]))
	return
}
