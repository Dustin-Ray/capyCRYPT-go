package main

import (
	"crypto/rand"
	"fmt"
)

// Generates n number of random bytes.
func GenerateRandomBytes(n uint) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	return b
}

// Converts byte array to uint64 array
func BytesToUint64(bytes []byte) (uint64s [25]uint64) {
	var result uint64
	for i := 0; i < len(bytes)/8; i += 8 {
		result = uint64(bytes[i]) | uint64(bytes[i+1])<<8 | uint64(bytes[i+2])<<16 |
			uint64(bytes[i+3])<<24 | uint64(bytes[i+4])<<32 | uint64(bytes[i+5])<<40 |
			uint64(bytes[i+6])<<48 | uint64(bytes[i+7])<<56
		uint64s[i] = result
	}
	return
}

func XORUint64Arrays(a, b []uint64) []uint64 {
	if len(a) != len(b) {
		return nil
	}

	result := make([]uint64, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] ^ b[i]
	}

	return result
}

func Uint64ToHex(input []uint64) string {
	output := ""
	for _, i := range input {
		output += fmt.Sprintf("%x", i)
	}
	return output
}
