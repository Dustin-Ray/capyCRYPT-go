package main

// Absorbs rate amount of data into the state. Returns
// pointer to state.
func SpongeAbsorb(m *[]byte, capacity int) *[25]uint64 {

	rateInBytes := (1600 - capacity) / 8
	P := *m
	if len(*m)%rateInBytes != 0 {
		P = padTenOne(*m, rateInBytes)
	}
	stateArray := BytesToStates(&P, rateInBytes)
	var S [25]uint64
	for _, st := range stateArray {
		S = Xorstates(S, st)
		KeccakF1600(&S)
	}
	return &S
}

// Squeezes bitlength amount of output from sponge with
// validity conditions 0 < bitLength < 2^2040
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

// Multi-rate padding scheme for sponge. Pads input
// to be multiple of rate.
func padTenOne(X []byte, rateInBytes int) []byte {
	q := rateInBytes - len(X)%rateInBytes
	padded := make([]byte, len(X)+q)
	copy(padded, X)
	padded[len(X)+q-1] = byte(0x80)
	return padded
}
