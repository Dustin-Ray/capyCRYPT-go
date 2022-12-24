package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"
)

func runtests() {

	testEncDec()

}

func testEncDec() {

	key := E521GenPoint(0)
	pw_string := []byte("test")
	pw_bytes := KMACXOF256(&pw_string, &[]byte{}, 512, "K")
	pw := big.NewInt(0)
	pw = pw.SetBytes(pw_bytes)
	pw = pw.Mul(pw, big.NewInt(4))
	pw = pw.Mod(pw, &key.n)

	key = key.SecMul(pw)
	message := []byte("test message")
	cgEnc := encryptKey(key, &message)
	// Create a buffer to write the data to
	var buf bytes.Buffer

	// Create a new encoder and use it to encode the data
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(cgEnc)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}

	data := buf.Bytes()
	fmt.Println("Array of bytes:", data)

	dec := gob.NewDecoder(&buf)
	var p2 Cryptogram
	if err := dec.Decode(&p2); err != nil {
		fmt.Println(err)
		return
	}

	// Get the encoded data as an array of bytes

	fmt.Println(p2)

	fmt.Println("working")
	for i := 0; i < 1; i++ {

		// fmt.Println(BytesToHexString(*test2.toBytes()))

		_, err := decryptKey(pw_string, &p2)
		if err != nil {
			fmt.Println("failed")
			break
		} else {
			fmt.Println("success")
		}
	}

}
