package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

func SHAKE(N *[]byte, d int) []byte {
	bytesToPad := 136 - len(*N)%136 // SHA3-256 r = 1088 / 8 = 136
	if bytesToPad == 1 {
		*N = append(*N, 0x86)
	} else {
		*N = append(*N, 0x06)
	}
	return SpongeSqueeze(SpongeAbsorb(N, 2*d), d, 1600-(2*d))
}

func SHAKE256(M []byte, d int) []byte {
	message := make([]byte, len(M)+1)
	copy(message, M)
	if (136 - len(M)%136) == 1 {
		message[len(M)] = 0x9f
	} else {
		message[len(M)] = 0x1f
	}
	return SpongeSqueeze(SpongeAbsorb(&message, d*2), d, 1600-(d*2))
}

func cSHAKE256(X []byte, L int, N string, S string) []byte {
	if N == "" && S == "" {
		return SHAKE256(X, L)
	}
	str := []byte(S)
	out := bytepad(append(encodeString([]byte(N)), encodeString(str)...), 136)
	out = append(out, X...)
	out = append(out, []byte{0x04}...) // https://keccak.team/keccak_specs_summary.html
	return SpongeSqueeze(SpongeAbsorb(&out, 512), L, 1600-512)
}

func KMACXOF256(K []byte, X []byte, L int, S string) []byte {
	newX := append(append(bytepad(encodeString(K), 136), X...), leftEncode(0)...)
	return cSHAKE256(newX, L, "KMAC", S)
}

func ComputeTaggedHash(pw, message []byte, S string) string {
	return hex.EncodeToString(KMACXOF256(pw, message, 512, S))
}

func ComputeSHA3HASH(data string, fileMode bool) string {
	dataBytes := []byte(data)
	if fileMode {
		return ""
	} else {
		return BytesToHexString(SHAKE(&dataBytes, 512))
	}
}

func encryptPW(pw string, message string) string {

	z := generateRandomBytes(64)
	K := []byte(pw)
	ke_ka := KMACXOF256(append(z, K...), []byte{}, 1024, "S")
	ke := ke_ka[:64]
	ka := ke_ka[64:]
	pW := KMACXOF256(ke, []byte(""), len([]byte(message))*8, "SKE")
	fmt.Println(hex.EncodeToString(pW))
	c := XorBytes(pW, []byte(message))
	fmt.Println(hex.EncodeToString(c))
	t := KMACXOF256(ka, []byte(message), 512, "SKA")

	return hex.EncodeToString(append(z, append(c, t...)...))
}

func decryptPW(pw string, message string) string {

	// var Replacer = strings.NewReplacer("\r\n", "")

	msg, _ := HexToBytes(message)

	K := []byte(pw)

	z := msg[:64]
	c := msg[64 : len(msg)-64]
	fmt.Println(hex.EncodeToString(c))
	t := msg[len(msg)-64:]

	ke_ka := KMACXOF256(append(z, K...), []byte{}, 1024, "S")
	ke := ke_ka[:64]
	ka := ke_ka[64:]
	pW := KMACXOF256(ke, []byte(""), len(c)*8, "SKE")
	fmt.Println(hex.EncodeToString(pW))
	m := XorBytes(pW, c)

	tP := KMACXOF256(ka, m, 512, "SKA")

	// fmt.Println("z: ", BytesToHexString(z))
	// fmt.Println("c: ", BytesToHexString(c))
	// fmt.Println("t: ", BytesToHexString(tP))

	if bytes.Equal(t, tP) {
		return string(m[:])
	} else {
		return string(m[:])
	}

}
