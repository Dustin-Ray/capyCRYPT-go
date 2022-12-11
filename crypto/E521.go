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

var (
	P, R, N big.Int
)

func initCurveConstants() {

	P := big.NewInt(2).Exp(big.NewInt(2), big.NewInt(521), nil)
	P.Sub(P, big.NewInt(1))

	R := big.NewInt(2).Exp(big.NewInt(2), big.NewInt(519), nil)
	s := big.NewInt(0)
	s.SetString("337554763258501705789107630418782636071904961214051226618635150085779108655765", 10)
	R.Sub(R, s)
	// N := R.Mul(R, big.NewInt(4))

}

// constructor for E521 with arbitrary return
func NewE521(x, y big.Int) *E521 {

	initCurveConstants()

	point := E521{
		X: x,
		Y: y,
		P: P,
		D: *big.NewInt(-376014),
		R: R,
	}
	return &point

}

// constructor for E521 with arbitrary return
func NewE521BasePoint(pX *big.Int, theLSB uint) *E521 {

	initCurveConstants()

	P := new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(521), nil), big.NewInt(1))

	X := pX
	num := new(big.Int).Sub(big.NewInt(1), new(big.Int).Exp(X, big.NewInt(2), nil))

	fmt.Println(num)
	fmt.Println(&P)

	num.Mod(num, P)

	denom := new(big.Int).Add(big.NewInt(376014), new(big.Int).Mul(X, new(big.Int).Exp(X, big.NewInt(2), nil)))
	denom.Mod(denom, P)

	denom = denom.ModInverse(denom, P)
	radicand := new(big.Int).Mul(num, denom)
	Y := sqrt(radicand, theLSB)

	point := E521{
		X: *big.NewInt(4),
		Y: *Y,
		P: *P,
		D: *big.NewInt(-376014),
		R: R,
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
	x2 := theOther.X

	y1 := e.Y
	y2 := theOther.Y

	xNum := (x1.Mul(&x1, &y2))
	xNum.Add(xNum, y1.Mul(&y1, &x2))
	xNum.Mod(xNum, &e.P)

	multiply := e.D.Mul(&e.D, &x1)
	multiply.Mul(multiply, &x2)
	multiply.Mul(multiply, &y1)

	xDenom := big.NewInt(1)
	xDenom.Add(xDenom, multiply.Mul(multiply, &y2))
	xDenom.Mod(xDenom, &e.P)
	xDenom.ModInverse(xDenom, &e.P)

	newX := xNum.Mul(xNum, xDenom)
	newX.Mod(newX, &e.P)

	yNum := y1.Mul(&y1, &y2)
	yNum.Sub(yNum, x1.Mul(&x1, &x2)).Mod(yNum, &e.P)

	yDenom := big.NewInt(1)
	yDenom.Sub(yDenom, multiply.Mul(multiply, &y2))
	yDenom.Mod(yDenom, &e.P)

	yDenom = yDenom.ModInverse(yDenom, &e.P)
	newY := yNum.Mul(yNum, yDenom).Mod(yNum, &e.P)

	return NewE521(*newX, *newY)
}

func MultiplyMontgomery(this *E521, theS *big.Int) *E521 {
	r0 := NewE521(*big.NewInt(0), *big.NewInt(1))
	r1 := this
	idx := theS.BitLen()
	for idx >= 0 {
		if theS.Bit(idx) == 0 {
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

func sqrt(v *big.Int, lsb uint) *big.Int {

	if v.Sign() == 0 {
		return big.NewInt(0)
	}
	r := new(big.Int).Exp(v, new(big.Int).Add(P.Rsh(&P, 2), big.NewInt(1)), &P)
	if r.Bit(0) != lsb {
		r.Sub(&P, r) // correct the lsb }

		bi := new(big.Int).Sub(new(big.Int).Mul(r, r), v)

		if bi.Mod(bi, &P).Sign() == 0 {
			return r
		} else {
			return nil
		}
	}
	return r
}

func main() {
	initCurveConstants()
	point := NewE521(*big.NewInt(4), *big.NewInt(4))
	fmt.Println(point.X)

}
