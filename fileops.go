package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/lukechampine/fastxor"
)

func generateRandomBytes(size int) []byte {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}

	return b
}

func StateArrayToHexString(input [25]uint64) string {
	var output string
	for _, v := range input {
		output += fmt.Sprintf("%x", v)
	}
	return output
}

func BytesToHexString(b []byte) string {
	return hex.EncodeToString(b)
}

func StateToByteArray(uint64s *[]uint64, bitLength int) []byte {
	var result []byte
	for _, v := range *uint64s {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, v)
		result = append(result, b...)
	}
	return result
}

func HexToBytes(hexString string) ([]byte, error) {
	return hex.DecodeString(hexString)
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
	for i := range a {
		result[i] ^= a[i] ^ b[i]
	}
	return result
}

func XorBytes(a, b []byte) []byte {
	dst := make([]byte, len(a))
	fastxor.Bytes(dst, a, b)
	return dst
}
