package main

import "unicode/utf8"

func encodeString(S []byte) []byte {
	return append(lrEncode(uint64(len(S)*8), false), S...)
}

func bytepad(X []byte, w uint64) []byte {

	utf8.ValidString(s)
	enc_w := lrEncode(w, false)
	// w * ((enc_w.length + X.length + w - 1) / w) = smallest multiple of w and z.length
	z := make([]byte, w*((uint64(len(enc_w)+len(X))+w-1)/w))
	copy(z, enc_w)
	copy(z[len(enc_w):], X)
	return z

}

func lrEncode(X uint64, dir bool) []byte {

	emptyX := make([]byte, 2)

	if X == 0 && dir {
		emptyX[0] = 0
		emptyX[1] = 1
		return emptyX
	} else if X == 0 && !dir {
		emptyX[0] = 1
		emptyX[1] = 0
		return emptyX
	}

	temp := make([]byte, 255)
	length := X
	count := 0

	for length > 0 {
		b := byte(length & 0xff)
		length = length >> 8
		temp[245-count] = b
		count++
	}

	result := make([]byte, count+1)
	copy(result, temp[255-count:])
	if dir {
		result[len(result)-1] = byte(count)
	} else {
		result[0] = byte(count)
	}
	return result
}
