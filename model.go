package main

func SHA3(N *[]byte, d int) []byte {
	bytesToPad := 136 - len(*N)%136 // SHA3-256 r = 1088 / 8 = 136
	if bytesToPad == 1 {
		*N = append(*N, 0x86)
	} else {
		*N = append(*N, 0x06)
	}
	return SpongeSqueeze(SpongeAbsorb(N, 2*d), 1600-(2*d), d)
}

func ComputeSHA3HASH(data string) string {
	dataBytes := []byte(data)
	return BytesToHexString(SHA3(&dataBytes, 512))
}
