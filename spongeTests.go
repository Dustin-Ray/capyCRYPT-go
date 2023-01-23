package main

import (
	"encoding/hex"
	"fmt"
)

func runSpongeTests() {

	str := "0xFFFFFFFFFFFFFF"
	temp, _ := hex.DecodeString(str)

	res := rightEncode(uint64(len(temp)))

	// hexstr := hex.EncodeToString(res)
	fmt.Println(res)
}
