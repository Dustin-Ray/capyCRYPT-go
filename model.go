package main

/**
Implements SHA3XOF functionality as defined in FIPS PUP 202 and NIST SP 800-105.
Inspiration:
https://github.com/mjosaarinen/tiny_sha3
https://keccak.team/keccak_specs_summary.html
https://github.com/NWc0de/KeccakUtils
Dustin Ray
version 0.1
*/

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"math/big"
)

type Cryptogram struct {
	Z_x *big.Int // Z_x is the x coordinate of the public nonce
	Z_y *big.Int // Z_y is the y coordinate of the public nonce
	Z   []byte   // optional Z public nonce for symmetric operations
	C   []byte   // c represents the ciphertext of an encryption
	T   []byte   // t is the authentication tag for the message
}

// Encodes a cryptogram to a byte array
func (cg *Cryptogram) encodeCryptogram() (*[]byte, error) {
	var buf bytes.Buffer
	// Create a new encoder and use it to encode the data
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(cg)
	if err != nil {
		return nil, errors.New("failed to encode cryptogram")
	}
	result := buf.Bytes()
	return &result, nil
}

// Parses a cryptogram from a string
func parseCryptogram(cg_dec *[]byte) (*Cryptogram, error) {
	buf := bytes.NewBuffer(*cg_dec)
	dec := gob.NewDecoder(buf)
	var p2 Cryptogram
	if err := dec.Decode(&p2); err != nil {
		return nil, errors.New("failed to decrypt")
	}
	return &p2, nil
}

/*
SHA3-Keccak functionaility ref NIST FIPS 202.

	N: pointer to message to be hashed.
	d: requested output length
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

/*
FIPS 202 Section 3 cSHAKE function returns customizable and
domain seperated length L SHA3XOF hash of input string.

	X: input message in bytes
	L: requested output length
	N: optional function name string
	S: option customization string
	return: SHA3XOF hash of length L of input message X
*/
func cSHAKE256(X *[]byte, L int, N string, S string) []byte {
	if N == "" && S == "" {
		return SHAKE(X, L)
	}
	out := bytepad(append(encodeString([]byte(N)), encodeString([]byte(S))...), 136)
	out = append(out, *X...)
	out = append(out, []byte{0x04}...) // https://keccak.team/keccak_specs_summary.html
	return SpongeSqueeze(SpongeAbsorb(&out, 512), L, 1600-512)
}

/*
Generates keyed hash for given input as specified in NIST SP 800-185 section 4.

	K: key
	X: byte-oriented message
	L: requested bit length
	S: customization string
	return: KMACXOF256 of X under K
*/
func KMACXOF256(K *[]byte, X *[]byte, L int, S string) []byte {
	newX := append(append(bytepad(encodeString(*K), 136), *X...), rightEncode(0)...)
	return cSHAKE256(&newX, L, "KMAC", S)
}

/*
Computes an authentication tag t of a byte array m under passphrase pw

	pw: symmetric encryption key, can be blank
	message: message to encrypt
	S: customization string
	return: t <- KMACXOF256(pw, m, 512, “T”)
*/
func ComputeTaggedHash(pw, message []byte, S string) string {
	return hex.EncodeToString(KMACXOF256(&pw, &message, 512, S))
}

/*
Computes SHA3-512 hash of byte array

	data: message to hash
	fileMode: determines wheter to process a file or text
	from the notepad.
	return: SHA3-512 hash of X
*/
func ComputeSHA3HASH(data string, fileMode bool) string {
	dataBytes := []byte(data)
	if fileMode {
		return ""
	} else {
		return hex.EncodeToString(SHAKE(&dataBytes, 512))
	}
}

/*
Encrypts a byte array m symmetrically under passphrase pw:

	z <- Random(512)
	(ke || ka) <- KMACXOF256(z || pw, “”, 1024, “S”)
	c <- KMACXOF256(ke, “”, |m|, “SKE”) xor m
	t <- KMACXOF256(ka, m, 512, “SKA”)
	pw: symmetric encryption key, can be blank
	message: message to encrypt
	return: symmetric cryptogram: (z, c, t)
*/
func encryptPW(pw []byte, msg *[]byte) *[]byte {

	z := generateRandomBytes(64)
	tempKeka := append(z, []byte(pw)...)
	ke_ka := KMACXOF256(&tempKeka, &[]byte{}, 1024, "S")
	ke := ke_ka[:64]
	ka := ke_ka[64:]
	pW := KMACXOF256(&ke, &[]byte{}, len(*msg)*8, "SKE")
	c := XorBytes(pW, *msg)
	t := KMACXOF256(&ka, msg, 512, "SKA")

	//construct a cryptogram
	result0 := Cryptogram{
		Z: z, C: c, T: t,
	}
	result, _ := result0.encodeCryptogram()
	return result
}

/*
Decrypts a symmetric cryptogram (z, c, t) under passphrase pw

	SECURITY NOTE: ciphertext length == plaintext length
	(ke || ka) <- KMACXOF256(z || pw, “”, 1024, “S”)
	m <- KMACXOF256(ke, “”, |c|, “SKE”) xor c
	t’ <- KMACXOF256(ka, m, 512, “SKA”)
	accept if, and only if, t’ = t
	msg: cryptogram to decrypt, assumes valid format.
	pw: decryption password, can be blank
	return: m, if and only if t` = t
*/
func decryptPW(pw []byte, cg *Cryptogram) (*[]byte, error) {

	z := cg.Z
	c := cg.C
	t := cg.T
	temp := append(z, pw...)
	ke_ka := KMACXOF256(&temp, &[]byte{}, 1024, "S")
	ke := ke_ka[:64]
	ka := ke_ka[64:]

	pW := KMACXOF256(&ke, &[]byte{}, len(c)*8, "SKE")
	m := XorBytes(c, pW)
	tP := KMACXOF256(&ka, &m, 512, "SKA")
	if bytes.Equal(t, tP) {
		return &m, nil
	} else {
		return &m, errors.New("unable to decrypt")
	}

}

/*
Generates a (Schnorr/ECDHIES) key pair from passphrase pw:

	s <- KMACXOF256(pw, “”, 512, “K”); s <- 4s
	V <- s*G

key pair: (s, V)
key: a pointer to an empty KeyObj to be populated with user data
*/
func generateKeyPair(ctx *WindowCtx, key *KeyObj) bool {
	return constructKey(ctx, key)
}

/*
Encrypts a byte array m under the (Schnorr/ECDHIES) public key V.
Operates under Schnorr/ECDHIES principle in that shared symmetric key is
exchanged with recipient. SECURITY NOTE: ciphertext length == plaintext length

	k <- Random(512); k <- 4k
	W <- k*V; Z <- k*G
	(ke || ka) <- KMACXOF256(W x , “”, 1024, “P”)
	c <- KMACXOF256(ke, “”, |m|, “PKE”) xor m
	t <- KMACXOF256(ka, m, 512, “PKA”)
	pubKey: X coordinate of public static key V, accepted as string
	message: message of any length or format to encrypt
	return cryptogram: (Z, c, t) = Z||c||t
*/
func encryptKey(pubKey *E521, message *[]byte) *[]byte {

	k := big.NewInt(0)
	k = k.SetBytes(generateRandomBytes(64))
	k = k.Mul(k, big.NewInt(4))
	k = k.Mod(k, &pubKey.n)

	W := pubKey.SecMul(k)

	Z := E521GenPoint(0).SecMul(k) //watch out for this, be sure to correct msb

	temp := W.x.Bytes()
	ke_ka := KMACXOF256(&temp, &[]byte{}, 1024, "P")
	ke := ke_ka[:64]
	ka := ke_ka[64:]

	c := XorBytes(KMACXOF256(&ke, &[]byte{}, len(*message)*8, "PKE"), *message)
	t := KMACXOF256(&ka, message, 512, "PKA")
	result0 := Cryptogram{Z_x: &Z.x, Z_y: &Z.y, C: c, T: t}
	result, _ := result0.encodeCryptogram()
	return result
}

/*
Decrypts a cryptogram under password. Assumes cryptogram is well-formed.
Operates under Schnorr/ECDHIES principle in that shared symmetric key is
derived from Z.

	s <- KMACXOF256(pw, “”, 512, “K”); s <- 4s
	W <- s*Z
	(ke || ka) <- KMACXOF256(W x , “”, 1024, “P”)
	m <- KMACXOF256(ke, “”, |c|, “PKE”) XOR c
	t’ <- KMACXOF256(ka, m, 512, “PKA”)
	@param pw password used to generate E521 encryption key.
	@param message cryptogram of format Z||c||t
	@return Decryption of cryptogram Z||c||t iff t` = t
*/
func decryptKey(pw []byte, message *Cryptogram) (*string, error) {

	Z := NewE521XY(*message.Z_x, *message.Z_y)

	s := big.NewInt(0)
	s = s.SetBytes(KMACXOF256(&pw, &[]byte{}, 512, "K"))
	s = s.Mul(s, big.NewInt(4))
	s = s.Mod(s, &Z.n)

	W := Z.SecMul(s)

	temp := W.x.Bytes()
	ke_ka := KMACXOF256(&temp, &[]byte{}, 1024, "P")
	ke := ke_ka[:64]
	ka := ke_ka[64:]
	m := XorBytes(KMACXOF256(&ke, &[]byte{}, len(message.C)*8, "PKE"), message.C)
	t_p := KMACXOF256(&ka, &m, 512, "PKA")
	if bytes.Equal(t_p, message.T) {
		result := string(m)
		return &result, nil
	} else {
		return nil, errors.New("decryption failure")
	}
}
