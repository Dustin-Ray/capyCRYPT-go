package main

import (
	"math/big"
)

/**
 * E521 Elliptic Curve (Edward's Curve) of equation: (x^2) + (y^2) = 1 + d(x^2)(y^2)
 * where d = -376014
 * Contains methods to add and multiply points on curve using scalar values.
 */
type E521 struct {
	x big.Int //X coordinate
	y big.Int // Y cooridinate
	p big.Int // Mersenne prime defining a finite field F(p)
	d big.Int // d = -376014
	r big.Int // number of points on Curve -> n := 4 * (R) .
	n big.Int //4 * r
}

// number of points on Curve -> n := 4 * (R) .
func (e *E521) getR() big.Int {
	R := big.NewInt(2).Exp(big.NewInt(2), big.NewInt(519), nil)
	s := big.NewInt(0)
	s.SetString("337554763258501705789107630418782636071904961214051226618635150085779108655765", 10)
	R.Sub(R, s)
	return *R
}

// Mersenne prime defining a finite field F(p)
func (e *E521) getP() big.Int {
	P := new(big.Int).Sub(big.NewInt(2).Exp(big.NewInt(2), big.NewInt(521), nil), big.NewInt(1))
	return *P
}

// constructor for E521 for any x, y
func NewE521XY(x, y big.Int) *E521 {
	tempR := new(E521).getR()
	point := E521{
		x: x,
		y: y,
		p: new(E521).getP(),
		d: *big.NewInt(-376014),
		r: tempR,
		n: *new(E521).r.Mul(&tempR, big.NewInt(4)),
	}
	return &point
}

// constructor for E521, solves for y
func NewE521X(x big.Int, msb uint) *E521 {
	tempR := new(E521).getR()
	point := E521{
		x: x,
		y: *solveForY(&x, new(E521).getP(), msb),
		p: new(E521).getP(),
		d: *big.NewInt(-376014),
		r: tempR,
		n: *new(E521).r.Mul(&tempR, big.NewInt(4)),
	}
	return &point
}

// Generator point for the curve, with x = 4 and y a unique even number obtained
// from solving curve equation.
func E521GenPoint(msb uint) *E521 {
	tempR := new(E521).getR()

	P := new(E521).getP()
	X := big.NewInt(4)
	Y := solveForY(X, P, msb)
	point := E521{
		x: *X,
		y: *Y,
		p: P,
		d: *big.NewInt(-376014),
		r: tempR,
		n: *new(E521).r.Mul(&tempR, big.NewInt(4)),
	}
	return &point

}

// solves curve equation for y value
func solveForY(X *big.Int, P big.Int, msb uint) *big.Int {
	num := new(big.Int).Sub(big.NewInt(1), new(big.Int).Exp(X, big.NewInt(2), nil))
	// fmt.Println("num: ", num)
	num = num.Mod(num, &P)
	// fmt.Println("num mod p: ", num)
	denom := new(big.Int).Add(big.NewInt(1), (new(big.Int).Mul(big.NewInt(376014), new(big.Int).Exp(X, big.NewInt(2), nil))))
	// fmt.Println("denom: ", denom)
	denom = denom.Mod(denom, &P)
	// fmt.Println("denom mod p: ", denom)
	denom = new(big.Int).ModInverse(denom, &P)
	// fmt.Println("denom mod inv: ", denom)
	radicand := new(big.Int).Mul(num, denom)
	// fmt.Println("radicand: ", radicand)
	Y := sqrt(radicand, msb)
	// fmt.Println("y: ", Y)
	return Y
}

// The identity point of the curve (also refered to as "point at infinity").
// Equivalent to 0 in integer group.
func E521IdPoint() *E521 { return NewE521XY(*big.NewInt(0), *big.NewInt(1)) }

/*
Gets the opposite value of a point, defined as the following:
if P = (X, Y), opposite of P = (-X, Y).
*/
func (e *E521) getOpposite() *E521 { return NewE521XY(*e.x.Neg(&e.x), e.y) }

// Checks two points for equality by comparing their coordinates.
func (A *E521) Equals(B *E521) bool { return A.x.Cmp(&B.x) == 0 && A.y.Cmp(&B.y) == 0 }

/*
Adds two E521 points and returns another E521 curve point.
Point addition operation is defined as:

	(x1, y1) + (x2, y2) = ((x1y2 + y1x2) / (1 + (d)x1x2y1y2)), ((y1y2 - x1x2) / (1 - (d)x1x2y1y2))

where "/" is defined to be multiplication by modular inverse.
*/
func (A *E521) Add(B *E521) *E521 {

	x1, y1, x2, y2 := A.x, A.y, B.x, B.y

	xNum := new(big.Int).Add(new(big.Int).Mul(&x1, &y2), new(big.Int).Mul(&y1, &x2))
	xNum.Mod(xNum, &A.p)

	mul := new(big.Int).Mul(&A.d, &x1) //x1 * x2 *  y1 * y2
	mul = new(big.Int).Mul(mul, &x2)
	mul = new(big.Int).Mul(mul, &y1)
	mul = new(big.Int).Mul(mul, &y2)

	xDenom := new(big.Int).Add(big.NewInt(1), mul)
	xDenom.Mod(xDenom, &A.p)
	xDenom = new(big.Int).ModInverse(xDenom, &A.p)

	newX := new(big.Int).Mul(xNum, xDenom)
	newX.Mod(newX, &A.p)

	yNum := new(big.Int).Sub(new(big.Int).Mul(&y1, &y2), new(big.Int).Mul(&x1, &x2))
	yNum.Mod(yNum, &A.p)

	yDenom := new(big.Int).Sub(big.NewInt(1), mul)
	yDenom.Mod(yDenom, &A.p)
	yDenom = new(big.Int).ModInverse(yDenom, &A.p)

	newY := new(big.Int).Mul(yNum, yDenom)
	newY.Mod(newY, &A.p)

	return NewE521XY(*newX, *newY)
}

/*
EC Multiplication algorithm using the Montgomery Ladder approach to mitigate
power consumption side channel attacks. Mostly constructed around:

(pg 4.)	https://eprint.iacr.org/2014/140.pdf

S is a  scalar value to multiply by. S is a private key and should be kept secret.
Returns Curve.E521 point which is result of multiplication.
*/
func (r1 *E521) SecMul(S *big.Int) *E521 {
	r0 := NewE521XY(*big.NewInt(0), *big.NewInt(1))
	for i := S.BitLen(); i >= 0; i-- {
		if S.Bit(i) == 1 {
			r0 = r0.Add(r1)
			r1 = r1.Add(r1)
		} else {
			r1 = r0.Add(r1)
			r0 = r0.Add(r0)
		}
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
	P := new(E521).getP()
	r := new(big.Int).Exp(v, new(big.Int).Add(new(big.Int).Rsh(&P, 2), big.NewInt(1)), &P)
	// fmt.Println("r value: ", r)
	if r.Bit(0) != lsb {
		r.Sub(&P, r) // correct the lsb }
		// fmt.Println("r sub value: ", r)
		bi := new(big.Int).Sub(new(big.Int).Mul(r, r), v)
		bi = bi.Mod(bi, &P)
		// fmt.Println("bi value: ", bi)
		if bi.Sign() == 0 {
			return r
		} else {
			return nil
		}
	}
	return r
}
