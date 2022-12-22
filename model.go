package main

/**
 * Implements SHA3XOF functionality as defined in FIPS PUP 202 and NIST SP 800-105.
 * Inspiration:
 *      <a href="https://github.com/mjosaarinen/tiny_sha3">...</a>
 *      <a href="https://keccak.team/keccak_specs_summary.html">...</a>
 *      <a href="https://github.com/NWc0de/KeccakUtils">...</a>
 * Dustin Ray
 * version 0.1
 */

import (
	"bytes"
	"encoding/hex"
	"errors"
)

/*
* SHA3-Keccak functionaility ref NIST FIPS 202.
* N: pointer to message to be hashed.
* d: requested output length
 */
func SHAKE(N *[]byte, d int) []byte {
	bytesToPad := 136 - len(*N)%136 // SHA3-256 r = 1088 / 8 = 136
	if bytesToPad == 1 {
		*N = append(*N, 0x86)
	} else {
		*N = append(*N, 0x06)
	}
	return SpongeSqueeze(SpongeAbsorb(N, 2*d), d, 1600-(2*d))
}

/**
 * FIPS 202 Section 3 cSHAKE function returns customizable and
 * domain seperated length L SHA3XOF hash of input string.
 * X: input message in bytes
 * L: requested output length
 * N: optional function name string
 * S: option customization string
 * return: SHA3XOF hash of length L of input message X
 */
func cSHAKE256(X []byte, L int, N string, S string) []byte {
	if N == "" && S == "" {
		return SHAKE(&X, L)
	}
	out := bytepad(append(encodeString([]byte(N)), encodeString([]byte(S))...), 136)
	out = append(out, X...)
	out = append(out, []byte{0x04}...) // https://keccak.team/keccak_specs_summary.html
	return SpongeSqueeze(SpongeAbsorb(&out, 512), L, 1600-512)
}

/**
 * Generates keyed hash for given input as specified in NIST SP 800-185 section 4.
 * K: key
 * X: byte-oriented message
 * L: requested bit length
 * S: customization string
 * return: KMACXOF256 of X under K
 */
func KMACXOF256(K []byte, X []byte, L int, S string) []byte {
	newX := append(append(bytepad(encodeString(K), 136), X...), rightEncode(0)...)
	return cSHAKE256(newX, L, "KMAC", S)
}

/**
 * Computes an authentication tag t of a byte array m under passphrase pw
 * pw: symmetric encryption key, can be blank
 * message: message to encrypt
 * S: customization string
 * return: t <- KMACXOF256(pw, m, 512, “T”)
 */
func ComputeTaggedHash(pw, message []byte, S string) string {
	return hex.EncodeToString(KMACXOF256(pw, message, 512, S))
}

/**
 * Computes SHA3-512 hash of byte array
 * data: message to hash
 * fileMode: determines wheter to process a file or text
 * from the notepad.
 * return: SHA3-512 hash of X
 */
func ComputeSHA3HASH(data string, fileMode bool) string {
	dataBytes := []byte(data)
	if fileMode {
		return ""
	} else {
		return hex.EncodeToString(SHAKE(&dataBytes, 512))
	}
}

/**
 * Encrypts a byte array m symmetrically under passphrase pw:
 * z <- Random(512)
 * (ke || ka) <- KMACXOF256(z || pw, “”, 1024, “S”)
 * c <- KMACXOF256(ke, “”, |m|, “SKE”) xor m
 * t <- KMACXOF256(ka, m, 512, “SKA”)
 * pw: symmetric encryption key, can be blank
 * message: message to encrypt
 * return: symmetric cryptogram: (z, c, t)
 */
func encryptPW(pw []byte, msg *[]byte) []byte {

	z := generateRandomBytes(64)
	ke_ka := KMACXOF256(append(z, []byte(pw)...), []byte{}, 1024, "S")
	ke := ke_ka[:64]
	ka := ke_ka[64:]
	pW := KMACXOF256(ke, []byte{}, len(*msg)*8, "SKE")
	c := XorBytes(pW, *msg)
	t := KMACXOF256(ka, *msg, 512, "SKA")
	result := append(z, append(c, t...)...)
	return result
}

/**
 * Decrypts a symmetric cryptogram (z, c, t) under passphrase pw
 * (ke || ka) <- KMACXOF256(z || pw, “”, 1024, “S”)
 * m <- KMACXOF256(ke, “”, |c|, “SKE”) xor c
 * t’ <- KMACXOF256(ka, m, 512, “SKA”)
 * accept if, and only if, t’ = t
 * message: cryptogram to decrypt, assumes valid format.
 * pw: decryption password, can be blank
 * return: m, if and only if t` = t
 */
func decryptPW(pw, msg []byte) ([]byte, error) {

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
		return m, nil
	} else {
		return nil, errors.New("decryption failure")
	}
}

/**
 * Generates a (Schnorr/ECDHIES) key pair from passphrase pw:
 *  s <- KMACXOF256(pw, “”, 512, “K”); s <- 4s
 *  V <- s*G
 *  key pair: (s, V)
 * key: a pointer to an empty KeyObj to be populated with user data
 */
func generateKeyPair(ctx *WindowCtx, key *KeyObj) bool {
	return constructKey(ctx, key)
}
