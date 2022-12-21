package main

func SpongeAbsorb(m *[]byte, capacity int) *[25]uint64 {

	rateInBytes := (1600 - capacity) / 8
	P := *m
	if len(*m)%rateInBytes != 0 {
		P = padTenOne(*m, rateInBytes)
	}
	stateArray := BytesToStates(P, rateInBytes)
	var S [25]uint64
	for _, st := range stateArray {
		S = Xorstates(S, st)
		KeccakF1600(&S)
	}
	return &S
}

func SpongeSqueeze(S *[25]uint64, bitLength, rate int) []byte {

	var out []uint64 //FIPS 202 Algorithm 8 Step 8
	offset := 0
	blockSize := rate / 64

	for len(out)*64 < bitLength {
		out = append(out, S[0:blockSize]...)
		offset += blockSize
		KeccakF1600(S) //FIPS 202 Algorithm 8 Step 10
	}
	return StateToByteArray(&out, bitLength/8)[:bitLength/8] //FIPS 202 3.1
}

func padTenOne(X []byte, rateInBytes int) []byte {
	q := rateInBytes - len(X)%rateInBytes
	padded := make([]byte, len(X)+q)
	copy(padded, X)
	padded[len(X)+q-1] = byte(0x80)
	return padded
}
