package main

import (
	"encoding/binary"
)

/*
NIST SP 800-185 2.3.2: String Encoding.
The encode_string function is used to encode bit strings in a way that may be parsed
unambiguously from the beginning of the string.

	return: left_encode(len(S)) + S.
*/
func encodeString(S []byte) []byte { return append(leftEncode(uint64(len(S)*8)), S...) }

/*
NIST SP 800-185 2.3.3
The bytepad(X, w) function prepends an encoding of the integer w to an input string X, then pads
the result with zeros until it is a byte string whose length in bytes is a multiple of w.

	X: the byte string to pad
	w: the rate of the sponge
	return: z = encode(X) + X + ("0" * LCM of length of z and w)
*/
func bytepad(input []byte, w int) []byte {
	// leftEncode always returns max 9 bytes
	buf := make([]byte, 0, 9+len(input)+w)
	buf = append(buf, leftEncode(uint64(w))...)
	buf = append(buf, input...)
	padlen := w - (len(buf) % w)
	return append(buf, make([]byte, padlen)...)
}

/*
leftEncode function is used to encode bit strings in a way that may be parsed
unambiguously from the beginning of the string by appending the encoding of
the length of the string to the end of the string.

	return: S + right_encode(len(S)).
*/
func rightEncode(value uint64) []byte {

	if value == 0 {
		return []byte{0, 1}
	}
	var b [9]byte
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

/*
leftEncode function is used to encode bit strings in a way that may be parsed
unambiguously from the beginning of the string by prepending the encoding of
the length of the string to the beginning of the string.

	return: left_encode(len(S)) + S.
*/
func leftEncode(value uint64) []byte {
	if value == 0 {
		return []byte{1, 0}
	}
	var b [9]byte
	binary.BigEndian.PutUint64(b[1:], value)
	// Trim all but last leading zero bytes
	i := byte(1)
	for i < 8 && b[i] == 0 {
		i++
	}
	// Append number of encoded bytes
	b[i-1] = 9 - i
	return b[i-1:]
}
