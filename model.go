package main

import (
	"bytes"
	"encoding/hex"
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
	out := bytepad(append(encodeString([]byte(N)), encodeString([]byte(S))...), 136)
	out = append(out, X...)
	out = append(out, []byte{0x04}...) // https://keccak.team/keccak_specs_summary.html
	return SpongeSqueeze(SpongeAbsorb(&out, 512), L, 1600-512)
}

func KMACXOF256(K []byte, X []byte, L int, S string) []byte {
	newX := append(append(bytepad(encodeString(K), 136), X...), rightEncode(0)...)
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

func encryptPW(pw, message []byte) []byte {

	z := generateRandomBytes(64)
	ke_ka := KMACXOF256(append(z, []byte(pw)...), []byte{}, 1024, "S")
	ke := ke_ka[:64]
	ka := ke_ka[64:]
	pW := KMACXOF256(ke, []byte{}, len(message)*8, "SKE")
	c := XorBytes(pW, message)
	t := KMACXOF256(ka, message, 512, "SKA")
	result := append(z, append(c, t...)...)
	return result
}

func decryptPW(pw, msg []byte) []byte {

	z := msg[:64]
	c := msg[64 : len(msg)-64]
	cP := make([]byte, len(msg)-128)
	copy(cP, c)
	t := msg[len(msg)-64:]
	ke_ka := KMACXOF256(append(z, pw...), []byte{}, 1024, "S")
	ke := ke_ka[:64]
	ka := ke_ka[64:]

	pW := KMACXOF256(ke, []byte{}, len(c)*8, "SKE")
	m := XorBytes(cP, pW)
	tP := KMACXOF256(ka, m, 512, "SKA")
	if bytes.Equal(t, tP) {
		return m
	} else {
		return m
	}
}

/**
 * Generates a (Schnorr/ECDHIES) key pair from passphrase pw:
 *  s <- KMACXOF256(pw, “”, 512, “K”); s <- 4s
 *  V <- s*G
 *  key pair: (s, V)
 */
func generateKeyPair(ctx *WindowCtx, key *KeyObj) bool {
	return constructKey(ctx.win, key)
}
