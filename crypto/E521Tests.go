package main

import (
	"fmt"
	"math/big"
)

func main() {

	Zero()
	One()
	GPlusMinusG()
	TwoTimesG()
	FourTimesG()
	NotZero()
	rTimesG()
	TestkTimesGAndkmodRTimesG()
	TestkPlus1TimesG()
	ktTimesgEqualskgtg()
	ktpEqualstkGEqualsktmodrG()

}

func Zero() {

	passedTestCount := 0
	numberOfTests := 100
	for i := 0; i < numberOfTests; i++ {
		G := IDPoint()
		if G.SecMul(big.NewInt(0)).Equals(IDPoint()) {
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
		G := GenPoint()
		if G.SecMul(big.NewInt(1)).Equals(GenPoint()) {
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
		G := GenPoint()
		if G.Add(GenPoint().getOpposite()).Equals(IDPoint()) {
			passedTestCount++
		} else {
			break
		}
	}
	fmt.Println("Test passed: ", passedTestCount == numberOfTests)
}

func TwoTimesG() {

	passedTestCount := 0
	numberOfTests := 100
	for i := 0; i < numberOfTests; i++ {
		G := GenPoint()
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
		G := GenPoint()
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
		G := GenPoint()
		if !G.SecMul(big.NewInt(4)).Equals(IDPoint()) {
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
		G := GenPoint()
		if G.SecMul(&G.r).Equals(IDPoint()) {
			passedTestCount++
		} else {
			break
		}
	}
	fmt.Println("Test passed: ", passedTestCount == numberOfTests)
}

func TestkTimesGAndkmodRTimesG() {
	G := GenPoint()
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
		G2 := GenPoint().SecMul(k)
		G2 = G2.Add(GenPoint())
		k = k.Add(k, big.NewInt(1))
		G1 := GenPoint().SecMul(k)
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

		G2 := GenPoint().SecMul(k)
		G2 = G2.Add(GenPoint().SecMul(t))

		x := new(big.Int).Add(k, t)
		G1 := GenPoint().SecMul(x)

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

		ktP := GenPoint().SecMul(t).SecMul(k)
		tkG := GenPoint().SecMul(k).SecMul(t)

		ktmodr := k.Mul(k, t)
		ktmodr = ktmodr.Mod(ktmodr, &GenPoint().r)
		ktmodrG := GenPoint().SecMul(ktmodr)

		if ktP.Equals(tkG) && ktP.Equals(ktmodrG) {
			passedTestCount++
		} else {
			break
		}
	}
	fmt.Println("Test passed: ", passedTestCount == numberOfTests)
}
