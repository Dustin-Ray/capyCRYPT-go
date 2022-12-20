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

func BytesToStates(in []byte, rateInBytes int) [][25]uint64 {
	stateArray := make([][25]uint64, (len(in) / rateInBytes)) //must accommodate enough states for datalength (in bytes) / rate
	offset := uint64(0)
	for i := 0; i < len(stateArray); i++ { //iterate through each state in stateArray
		var state [25]uint64                      // init empty state
		for j := 0; j < (rateInBytes*8)/64; j++ { //fill each state with rate # of bits
			state[j] = BytesToLane(in, offset)
			offset += 8
		}
		stateArray[i] = state
	}
	return stateArray
}

func BytesToLane(in []byte, offset uint64) uint64 {
	lane := uint64(0)
	for i := uint64(0); i < uint64(8); i++ {
		lane += uint64(in[i+offset]&0xFF) << (8 * i) //mask shifted byte to long and add to lane
	}
	return lane
}

func Xorstates(a, b [25]uint64) [25]uint64 {
	var result [25]uint64
	for i := range result {
		result[i] ^= a[i] ^ b[i]
	}
	return result
}

func padTenOne(X []byte, rateInBytes int) []byte {
	q := rateInBytes - len(X)%rateInBytes
	padded := make([]byte, len(X)+q)
	copy(padded, X)
	padded[len(X)+q-1] = byte(0x80)
	return padded
}
