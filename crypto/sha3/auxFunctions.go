package CryptoTool

import (
	"crypto/rand"
	"fmt"
)

func generateRandomBytes() []byte {
	b := make([]byte, 1000)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	return b
}
