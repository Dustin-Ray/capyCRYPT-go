package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func rune521Tests() {

	// fmt.Println("value for r:", E521IdPoint().r.String())
	// fmt.Println("value for p:", E521IdPoint().p.String())
	// fmt.Println("value for n:", E521IdPoint().n.String())
	// point := E521GenPoint(0)
	// fmt.Println("E521 gen point: ", point.y.String())

	// neg_g := E521GenPoint(0).getOpposite()

	// res := point.Add(neg_g)

	// Zero()
	// One()
	// GPlusMinusG()
	// TwoTimesG()
	// FourTimesG()
	// NotZero()
	// rTimesG()
	// TestkTimesGAndkmodRTimesG()
	// TestkPlus1TimesG()
	// ktTimesgEqualskgtg()
	// ktpEqualstkGEqualsktmodrG()

}

func Zero() {

	passedTestCount := 0
	numberOfTests := 100
	for i := 0; i < numberOfTests; i++ {
		G := E521IdPoint()
		if G.SecMul(big.NewInt(0)).Equals(E521IdPoint()) {
			passedTestCount++
		} else {
			break
		}
	}
	fmt.Println("Test passed: ", passedTestCount == numberOfTests)
}

func One() {

	passedTestCount := 0
	numberOfTests := 100
	for i := 0; i < numberOfTests; i++ {
		G := E521GenPoint(0)
		if G.SecMul(big.NewInt(1)).Equals(E521GenPoint(0)) {
			passedTestCount++
		} else {
			break
		}
	}
	fmt.Println("Test passed: ", passedTestCount == numberOfTests)
}

func GPlusMinusG() {

	passedTestCount := 0
	numberOfTests := 100
	for i := 0; i < numberOfTests; i++ {
		G := E521GenPoint(0)
		if G.Add(E521GenPoint(0).getOpposite()).Equals(E521IdPoint()) {
			passedTestCount++
		} else {
			break
		}
	}
	fmt.Println("Test passed: ", passedTestCount == numberOfTests)
}

func TwoTimesG() {

	passedTestCount := 0
	numberOfTests := 1
	for i := 0; i < numberOfTests; i++ {
		G := E521GenPoint(0)
		p := G.SecMul(big.NewInt(2))
		fmt.Println(p.x.String())
		fmt.Println(p.y.String())
		if G.SecMul(big.NewInt(2)).Equals(G.Add(G)) {
			passedTestCount++
		} else {
			break
		}
	}
	fmt.Println("Test passed: ", passedTestCount == numberOfTests)
}

func FourTimesG() {

	passedTestCount := 0
	numberOfTests := 100
	for i := 0; i < numberOfTests; i++ {
		G := E521GenPoint(0)
		if G.SecMul(big.NewInt(4)).Equals(G.SecMul(big.NewInt(2)).SecMul(big.NewInt(2))) {
			passedTestCount++
		} else {
			break
		}
	}
	fmt.Println("Test passed: ", passedTestCount == numberOfTests)
}

func NotZero() {

	passedTestCount := 0
	numberOfTests := 100
	for i := 0; i < numberOfTests; i++ {
		G := E521GenPoint(0)
		if !G.SecMul(big.NewInt(4)).Equals(E521IdPoint()) {
			passedTestCount++
		} else {
			break
		}
	}
	fmt.Println("Test passed: ", passedTestCount == numberOfTests)

}

func rTimesG() {

	passedTestCount := 0
	numberOfTests := 100
	for i := 0; i < numberOfTests; i++ {
		G := E521GenPoint(0)
		if G.SecMul(&G.r).Equals(E521IdPoint()) {
			passedTestCount++
		} else {
			break
		}
	}
	fmt.Println("Test passed: ", passedTestCount == numberOfTests)
}

func TestkTimesGAndkmodRTimesG() {
	G := E521GenPoint(0)
	R := G.getR()

	passedTestCount := 0
	numberOfTests := 50
	for i := 0; i < numberOfTests; i++ {
		k := generateRandomBigInt()
		G1 := G.SecMul(k)
		G2 := G.SecMul(k.Mod(k, &R))
		if G1.Equals(G2) {
			passedTestCount++
		} else {
			break
		}
	}
	fmt.Println("Test passed: ", passedTestCount == numberOfTests)
}

func TestkPlus1TimesG() {

	passedTestCount := 0
	numberOfTests := 50
	for i := 0; i < numberOfTests; i++ {
		k := generateRandomBigInt()
		G2 := E521GenPoint(0).SecMul(k)
		G2 = G2.Add(E521GenPoint(0))
		k = k.Add(k, big.NewInt(1))
		G1 := E521GenPoint(0).SecMul(k)
		if G1.Equals(G2) {
			passedTestCount++
		} else {
			break
		}
	}
	fmt.Println("Test passed: ", passedTestCount == numberOfTests)
}

func ktTimesgEqualskgtg() {

	passedTestCount := 0
	numberOfTests := 50
	for i := 0; i < numberOfTests; i++ {
		k := generateRandomBigInt()
		t := generateRandomBigInt()

		G2 := E521GenPoint(0).SecMul(k)
		G2 = G2.Add(E521GenPoint(0).SecMul(t))

		x := new(big.Int).Add(k, t)
		G1 := E521GenPoint(0).SecMul(x)

		if G1.Equals(G2) {
			passedTestCount++
		} else {
			break
		}
	}
	fmt.Println("Test passed: ", passedTestCount == numberOfTests)
}

func ktpEqualstkGEqualsktmodrG() {

	passedTestCount := 0
	numberOfTests := 50
	for i := 0; i < numberOfTests; i++ {
		k := generateRandomBigInt()
		t := generateRandomBigInt()

		ktP := E521GenPoint(0).SecMul(t).SecMul(k)
		tkG := E521GenPoint(0).SecMul(k).SecMul(t)

		ktmodr := k.Mul(k, t)
		ktmodr = ktmodr.Mod(ktmodr, &E521GenPoint(0).r)
		ktmodrG := E521GenPoint(0).SecMul(ktmodr)

		if ktP.Equals(tkG) && ktP.Equals(ktmodrG) {
			passedTestCount++
		} else {
			break
		}
	}
	fmt.Println("Test passed: ", passedTestCount == numberOfTests)
}

// gengerates random 512 bit integer
func generateRandomBigInt() *big.Int {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	random := big.NewInt(0)
	random.SetBytes(b)
	return random
}
