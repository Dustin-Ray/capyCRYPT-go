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
		G := IdentityPoint()
		if G.MultiplyMontgomery(big.NewInt(0)).Equals(IdentityPoint()) {
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
		if G.MultiplyMontgomery(big.NewInt(1)).Equals(GenPoint()) {
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
		if G.Add(GenPoint().getOpposite()).Equals(IdentityPoint()) {
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
		if G.MultiplyMontgomery(big.NewInt(2)).Equals(G.Add(G)) {
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
		if G.MultiplyMontgomery(big.NewInt(4)).Equals(G.MultiplyMontgomery(big.NewInt(2)).MultiplyMontgomery(big.NewInt(2))) {
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
		if !G.MultiplyMontgomery(big.NewInt(4)).Equals(IdentityPoint()) {
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
		if G.MultiplyMontgomery(&G.r).Equals(IdentityPoint()) {
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
		G1 := G.MultiplyMontgomery(k)
		G2 := G.MultiplyMontgomery(k.Mod(k, &R))
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
		G2 := GenPoint().MultiplyMontgomery(k)
		G2 = G2.Add(GenPoint())
		k = k.Add(k, big.NewInt(1))
		G1 := GenPoint().MultiplyMontgomery(k)
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

		G2 := GenPoint().MultiplyMontgomery(k)
		G2 = G2.Add(GenPoint().MultiplyMontgomery(t))

		x := new(big.Int).Add(k, t)
		G1 := GenPoint().MultiplyMontgomery(x)

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

		ktP := GenPoint().MultiplyMontgomery(t).MultiplyMontgomery(k)
		tkG := GenPoint().MultiplyMontgomery(k).MultiplyMontgomery(t)

		ktmodr := k.Mul(k, t)
		ktmodr = ktmodr.Mod(ktmodr, &GenPoint().r)
		ktmodrG := GenPoint().MultiplyMontgomery(ktmodr)

		if ktP.Equals(tkG) && ktP.Equals(ktmodrG) {
			passedTestCount++
		} else {
			break
		}
	}
	fmt.Println("Test passed: ", passedTestCount == numberOfTests)

}
