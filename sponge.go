package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func SHA3(N *[]byte, d int) []byte {
	bytesToPad := 136 - len(*N)%136 // SHA3-256 r = 1088 / 8 = 136
	if bytesToPad == 1 {
		*N = append(*N, 0x86)
	} else {
		*N = append(*N, 0x06)
	}
	return SpongeSqueeze(SpongeAbsorb(N, 2*d), d, 1600-(2*d))
}

func SpongeAbsorb(m *[]byte, capacity int) [25]uint64 {

	rateInBytes := (1600 - capacity) / 8
	var P []byte
	if len(*m)%rateInBytes == 0 {
		P = *m
	} else {
		P = padTenOne(*m, rateInBytes)
	}
	stateArray := BytesToStates(P, rateInBytes)
	var S [25]uint64
	for _, st := range stateArray {
		S = Xorstates(S, st)
		KeccakF1600(&S)
	}
	return S
}

func SpongeSqueeze(S [25]uint64, rate, bitLength int) []byte {

	var out []uint64 //FIPS 202 Algorithm 8 Step 8
	offset := 0
	blockSize := rate / 64

	for len(out)*64 < bitLength {
		out = append(out, S[0:blockSize]...)
		offset += blockSize
		KeccakF1600(&S) //FIPS 202 Algorithm 8 Step 10
	}
	return StateToByteArray(out, bitLength/8)[:512/8] //FIPS 202 3.1
}

func StateToByteArray(uint64s []uint64, bitLength int) []byte {
	var result []byte
	for _, v := range uint64s {
		// Use binary.PutUvarint to convert uint64 to byte array
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, v)
		result = append(result, b...)
	}
	return result
}

func BytesToStates(in []byte, rateInBytes int) [][25]uint64 {
	stateArray := make([][25]uint64, (len(in) / rateInBytes)) //must accommodate enough states for datalength (in bytes) / rate
	offset := uint64(0)
	for i := 0; i < len(stateArray); i++ { //iterate through each state in stateArray
		var state [25]uint64                      // init empty state
		for j := 0; j < (rateInBytes*8)/64; j++ { //fill each state with rate # of bits
			state[j] = BytesToLane(in, offset)
			offset += 8
		}
		stateArray[i] = state
	}

	return stateArray
}

func BytesToLane(in []byte, offset uint64) uint64 {
	lane := uint64(0)
	for i := uint64(0); i < uint64(8); i++ {
		lane += uint64(in[i+offset]&0xFF) << (8 * i) //mask shifted byte to long and add to lane
	}
	return lane
}

func Xorstates(a, b [25]uint64) [25]uint64 {
	var result [25]uint64
	for i := range result {
		result[i] ^= a[i] ^ b[i]
	}
	return result
}

func padTenOne(X []byte, rateInBytes int) []byte {
	q := rateInBytes - len(X)%rateInBytes
	padded := make([]byte, len(X)+q)
	copy(padded, X)
	padded[len(X)+q-1] = byte(0x80)
	return padded
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

func encode_string(S []byte) []byte { return append(lrEncode(uint64(len(S)*8), false), S...) }

func bytepad(X []byte, w uint64) []byte {

	enc_w := lrEncode(w, false)
	// w * ((enc_w.length + X.length + w - 1) / w) = smallest multiple of w and z.length
	z := make([]byte, w*((uint64(len(enc_w)+len(X))+w-1)/w))
	copy(z, enc_w)
	copy(z[len(enc_w):], X)
	return z

}

func ArrayToHexString(input [25]uint64) string {
	var output string
	for _, v := range input {
		output += fmt.Sprintf("%x", v)
	}
	return output
}

func generateRandomBytes() []byte {
	b := make([]byte, 136)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	return b
}

func BytesToHexString(b []byte) string {
	return hex.EncodeToString(b)
}
