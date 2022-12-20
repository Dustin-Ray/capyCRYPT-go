package main

import "encoding/hex"

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
	newX := append(append(bytepad(encodeString(K), 136), X...), lrEncode(0, true)...)
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
