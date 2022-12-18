package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func generateRandomBytes() []byte {
	b := make([]byte, 136)
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
