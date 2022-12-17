package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func SpongeAbsorb(m *[]byte, capacity int) *[25]uint64 {

	rateInBytes := (1600 - capacity) / 8
	var P *[]byte
	if len(*m)%rateInBytes == 0 {
		*P = *m
	} else {
		*P = padTenOne(*m, rateInBytes)
	}
	// fmt.Println(BytesToHexString(P))

	stateArray := BytesToStates(*P, rateInBytes)
	var S *[25]uint64
	// fmt.Println("stateArray:\n", stateArray)
	for _, st := range stateArray {
		keccakf(Xorstates(S, &st))
	}
	// fmt.Println(S)
	return S
}

func SpongeSqueeze(S *[25]uint64, rate, bitlength int) []byte {
	output := make([]byte, bitlength/8)

	for i := 0; i < len(output); i += 8 {
		block := S[i/8 : i/8+8]
		for j := 0; j < len(block); j++ {
			output[i+j] = byte(block[j] >> (uint(j) * 8))
		}
		if i+8 >= bitlength {
			break
		}
		keccakf(S)
	}
	// fmt.Println(output)
	return output

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

func SHA3(N []byte, d int) []byte {
	message := make([]byte, len(N)+1)
	copy(message, N)
	bytesToPad := 136 - len(N)%136 // SHA3-256 r = 1088 / 8 = 136
	if bytesToPad == 1 {
		message[len(N)] = 0x86
	} else {
		message[len(N)] = 0x06
	}
	return SpongeSqueeze(SpongeAbsorb(&message, 2*d), d, 1600-(2*d))
}

func BytesToLane(in []byte, offset uint64) uint64 {
	lane := uint64(0)
	for i := uint64(0); i < uint64(8); i++ {
		lane += uint64(in[i+offset]&0xFF) << (8 * i) //mask shifted byte to long and add to lane
		// fmt.Println(lane)
	}
	return lane
}

func Xorstates(a, b *[25]uint64) *[25]uint64 {
	var result [25]uint64
	for i := range result {
		result[i] ^= a[i] ^ b[i]
	}
	return &result
}

func keccakf(state *[25]uint64) {

}

func padTenOne(X []byte, rateInBytes int) []byte {
	q := rateInBytes - len(X)%rateInBytes
	padded := make([]byte, len(X)+q)
	copy(padded, X)
	b := byte(0x80)
	padded[len(X)+q-1] = b
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

func main() {

	msg := []byte{}
	// Keccak(msg, 1344)
	secParam := 512
	// fmt.Println(BytesToHexString(SpongeSqueeze(SpongeAbsorb(msg, 512), 1600-secParam, 256)))
	fmt.Println(SHA3(msg, secParam))

	var x uint64 = 8796093022209
	var y int64 = 8796093022215

	fmt.Printf("uint64: %v = %#[1]x, int64: %v = %#x\n", x, y, uint64(y))

}

func ArrayToHexString(input [25]uint64) string {
	var output string
	for _, v := range input {
		output += fmt.Sprintf("%x", v)
	}
	return output
}

func Uint64ToBytes(uint64s [25]uint64) []byte {
	var bytes []byte
	for i := range uint64s {
		bytes = append(bytes, byte(uint64s[i]))
	}
	return bytes
}

func generateRandomBytes() []byte {
	b := make([]byte, 136)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	// fmt.Println("random bytes: ", BytesToHexString(b))
	return b
}

func BytesToHexString(b []byte) string {
	return hex.EncodeToString(b)
}
