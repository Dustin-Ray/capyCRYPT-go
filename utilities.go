package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/lukechampine/fastxor"
)

// generates size number of random bytes. Is not assumed
// to be cryptographically secure.
func generateRandomBytes(size int) []byte {
	b := make([]byte, size)
	rand.Read(b)
	return b
}

// Converts uint64 arrays to hex strings
func StateArrayToHexString(input [25]uint64) string {
	var output string
	for _, v := range input {
		output += fmt.Sprintf("%x", v)
	}
	return output
}

// Encodes bytes to hex characters.
func BytesToHexString(b []byte) string {
	return hex.EncodeToString(b)
}

// Converts a single state to an array of bytes
func StateToByteArray(uint64s *[]uint64, bitLength int) []byte {
	var result []byte
	for _, v := range *uint64s {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, v)
		result = append(result, b...)
	}
	return result
}

// Converts a string of hex characters to a byte array.
func HexToBytes(hexString string) []byte {
	result, _ := hex.DecodeString(hexString)
	return result
}

// Main entry point for file and text processing. Converts byte array to
// series of state arrays per FIPS 202 format.
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

// Converts bytes to 64 bit lane/word
func BytesToLane(in []byte, offset uint64) uint64 {
	lane := uint64(0)
	for i := uint64(0); i < uint64(8); i++ {
		lane += uint64(in[i+offset]&0xFF) << (8 * i) //mask shifted byte to long and add to lane
	}
	return lane
}

// returns a XOR b, assumes equal size
func Xorstates(a, b [25]uint64) [25]uint64 {

	var result [25]uint64
	for i := range a {
		result[i] ^= a[i] ^ b[i]
	}
	return result
}

// returns a fastxor b as byte array. Assumes equal length.
func XorBytes(a, b []byte) []byte {
	dst := make([]byte, len(a))
	fastxor.Bytes(dst, a, b)
	return dst
}
