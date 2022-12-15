package main

import (
	"crypto/rand"
	"fmt"
)

// generates 1000 random bytes
func GenerateRandomBytes() []byte {
	b := make([]byte, 1000)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	return b
}
