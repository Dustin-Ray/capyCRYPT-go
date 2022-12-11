package main

import (
	"fmt"
	"math/big"
)

/**
 * E521 Elliptic Curve (Edward's Curve) of equation: (x^2) + (y^2) = 1 + d(x^2)(y^2)
 * where d = -376014
 * Contains methods to add and multiply points on curve using scalar values.
 */

type E521 struct {
	X big.Int //X coordinate
	Y big.Int // Y cooridinate
	P big.Int // Mersenne prime defining a finite field F(p)
	D big.Int // d = -376014
	R big.Int // number of points on Curve -> n := 4 * (R) .
}

func getR() big.Int {

	R := big.NewInt(2).Exp(big.NewInt(2), big.NewInt(519), nil)
	s := big.NewInt(0)
	s.SetString("337554763258501705789107630418782636071904961214051226618635150085779108655765", 10)
	R.Sub(R, s)
	return *R

}

func getP() big.Int {
	P := new(big.Int).Sub(big.NewInt(2).Exp(big.NewInt(2), big.NewInt(521), nil), big.NewInt(1))
	return *P
}

// constructor for E521 with arbitrary return
func NewE521(x, y big.Int) *E521 {

	point := E521{
		X: x,
		Y: y,
		P: getP(),
		D: *big.NewInt(-376014),
		R: getR(),
	}
	return &point

}

// constructor for generator point for E521 curve
func NewE521BasePoint(pX *big.Int, theLSB uint) *E521 {

	P := new(big.Int).Sub(big.NewInt(2).Exp(big.NewInt(2), big.NewInt(521), nil), big.NewInt(1))

	X := pX
	num := new(big.Int).Sub(big.NewInt(1), new(big.Int).Exp(X, big.NewInt(2), nil))
	num = num.Mod(num, P)

	denom := new(big.Int).Add(big.NewInt(1), (new(big.Int).Mul(big.NewInt(376014), new(big.Int).Exp(X, big.NewInt(2), nil))))
	denom = denom.Mod(denom, P)
	denom = new(big.Int).ModInverse(denom, P)
	radicand := new(big.Int).Mul(num, denom)

	Y := sqrt(radicand, theLSB)
	point := E521{
		X: *big.NewInt(4),
		Y: *Y,
		P: getP(),
		D: *big.NewInt(-376014),
		R: getR(),
	}
	return &point

}

func (e *E521) getOpposite() *E521 { return NewE521(*e.X.Neg(&e.X), e.Y) }

/*
 * Adds two E521 points and returns another (X, Y) point which is on the Curve.E521 curve defined
 * using parameters set by constructor. Add operation is defined as:
 * (x1, y1) + (x2, y2) = ((x1y2 + y1x2) / (1 + (d)x1x2y1y2)), ((y1y2 - x1x2) / (1 - (d)x1x2y1y2))
 * where "/" is defined to be multiplication by modular inverse.
 */
func (e *E521) Add(theOther *E521) *E521 {

	x1 := e.X
	fmt.Println("x1: " + x1.String())
	x2 := theOther.X
	fmt.Println("x2: " + x2.String())

	y1 := e.Y
	fmt.Println("y1: " + y1.String())
	y2 := theOther.Y
	fmt.Println("y2: " + y2.String())

	xNum := new(big.Int).Add(new(big.Int).Mul(&x1, &y2), new(big.Int).Mul(&y1, &x2))
	xNum.Mod(xNum, &e.P)

	fmt.Println("xNum: " + xNum.String())

	multiply := new(big.Int).Mul(&e.D, &x1)
	new(big.Int).Mul(multiply, &x2)
	new(big.Int).Mul(multiply, &y1)

	xDenom := new(big.Int).Add(big.NewInt(1), new(big.Int).Mul(multiply, &y2))
	xDenom.Mod(xDenom, &e.P)

	fmt.Println("xDenom : " + xDenom.String())

	xDenom = xDenom.ModInverse(xDenom, &e.P)

	fmt.Println("xDenom mod inverse: " + xDenom.String())

	newX := new(big.Int).Mul(xNum, xDenom)
	newX.Mod(newX, &e.P)

	fmt.Println("newX: " + newX.String())

	return NewE521(*newX, *big.NewInt(0))
}

/*
 * EC Multiplication algorithm using the Montgomery Ladder approach to mitigate
 * power consumption side channel attacks. Mostly constructed around
 * <a href="https://eprint.iacr.org/2014/140.pdf">...</a> pg 4.
 * theS is a  scalar value to multiply by. S is a private key and should be kept secret.
 * return Curve.E521 point which is result of multiplication.
 */
func (r1 *E521) MultiplyMontgomery(S *big.Int) *E521 {
	r0 := NewE521(*big.NewInt(0), *big.NewInt(1))
	idx := S.BitLen()
	for idx >= 0 {
		if S.Bit(idx) == 0 {
			r0 = r0.Add(r1)
			r1 = r1.Add(r1)
		} else {
			r1 = r0.Add(r1)
			r0 = r0.Add(r0)
		}
		idx--
	}
	return r0 // r0 = P * s
}

/*
 * Compute a square root of v mod p with a specified
 * the least significant bit, if such a root exists.
 * Provided by Dr. Paulo Barretto.
 * @param v   the radicand.
 * lsb is desired least significant bit (true: 1, false: 0).
 * return a square root r of v mod p with r mod 2 = 1 iff lsb = true
 * if such a root exists, otherwise null.
 */
func sqrt(v *big.Int, lsb uint) *big.Int {

	if v.Sign() == 0 {
		return big.NewInt(0)
	}

	P := new(big.Int).Sub(big.NewInt(2).Exp(big.NewInt(2), big.NewInt(521), nil), big.NewInt(1))

	r := new(big.Int).Exp(v, new(big.Int).Add(new(big.Int).Rsh(P, 2), big.NewInt(1)), P)
	if r.Bit(0) != lsb {
		r.Sub(P, r) // correct the lsb }

		bi := new(big.Int).Sub(new(big.Int).Mul(r, r), v)

		if bi.Mod(bi, P).Sign() == 0 {
			return r
		} else {
			return nil
		}
	}
	return r
}

func main() {

	point := NewE521BasePoint(big.NewInt(4), 0)
	point = point.MultiplyMontgomery(big.NewInt(1))

}
