package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"
)

func runtests() {

	testSig()

}

func testSig() {

	message := []byte("test")
	pw := []byte("password")

	s := new(big.Int).SetBytes(KMACXOF256(&pw, &[]byte{}, 512, "K"))
	s = s.Mul(s, big.NewInt(4))
	s = s.Mod(s, &E521IdPoint().n)

	V := *E521GenPoint(0)
	V = *V.SecMul(s)

	//get signing key for messsage under password
	k := new(big.Int).SetBytes(KMACXOF256(&pw, &message, 512, "N"))
	k = new(big.Int).Mul(k, big.NewInt(4))
	k = new(big.Int).Mod(k, &E521IdPoint().n)

	//create public signing key for message
	U := E521GenPoint(0).SecMul(k)
	fmt.Println(U.x)
	uXBytes := U.x.Bytes()

	//get the tag for the message key
	h := KMACXOF256(&uXBytes, &message, 512, "T")

	//create public nonce for signature
	h_bigInt := new(big.Int).SetBytes(h)
	z := new(big.Int).Sub(h_bigInt, k.Mul(k, s))
	z = new(big.Int).Mod(z, &E521IdPoint().r)

	// z = (k - hs) mod r
	result0 := Signature{H: new(big.Int).Abs(h_bigInt), Z: z}
	// result, err := encodeSignature(&result0)

	U2 := E521GenPoint(0).SecMul(result0.Z).Add(V.SecMul(result0.H))
	keyBytes := U2.x.Bytes()
	h_p := KMACXOF256(&keyBytes, &message, 512, "T")
	h2 := new(big.Int).SetBytes(h_p)
	h2 = new(big.Int).Abs(h2)
	if h2 == new(big.Int).SetBytes(h_p) {
		fmt.Println("success")
	}

	// if err != nil {
	// 	fmt.Println("failed to encode signature")
	// } else {
	// 	// fmt.Println(result)
	// }

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
	cgEnc := encryptWithKey(key, &message)
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

		_, err := encryptWithPassword(pw_string, &p2)
		if err != nil {
			fmt.Println("failed")
			break
		} else {
			fmt.Println("success")
		}
	}

}
