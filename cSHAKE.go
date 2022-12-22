package main

import (
	"encoding/binary"
)

func encodeString(S []byte) []byte { return append(leftEncode(uint64(len(S)*8)), S...) }

func bytepad(input []byte, w int) []byte {
	// leftEncode always returns max 9 bytes
	buf := make([]byte, 0, 9+len(input)+w)
	buf = append(buf, leftEncode(uint64(w))...)
	buf = append(buf, input...)
	padlen := w - (len(buf) % w)
	return append(buf, make([]byte, padlen)...)
}

func rightEncode(value uint64) []byte {

	if value == 0 {
		return []byte{0, 1}
	}

	var b [9]byte
	// big endian -> 00001
	// le -> 10000
	binary.BigEndian.PutUint64(b[0:], value)
	// Trim all but last leading zero bytes
	i := byte(1)
	for i < 8 && b[i] == 0 {
		i++
	}
	// Prepend number of encoded bytes
	b[0] = 9 - i
	return b[:9-i]
}

func leftEncode(value uint64) []byte {
	var b [9]byte
	binary.BigEndian.PutUint64(b[1:], value)
	// Trim all but last leading zero bytes
	i := byte(1)
	for i < 8 && b[i] == 0 {
		i++
	}
	// Prepend number of encoded bytes
	b[i-1] = 9 - i
	return b[i-1:]
}
