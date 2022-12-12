package main

import (
	"crypto/rand"
	"fmt"
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
	s string  //string representation of point
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

// constructor for E521 with arbitrary return
func NewE521(x, y big.Int) *E521 {
	point := E521{
		x: x,
		y: y,
		p: new(E521).getP(),
		d: *big.NewInt(-376014),
		r: new(E521).getR(),
		s: "x coord:\n" + x.String() + "\ny coord:\n" + y.String(),
	}
	return &point
}

// Generator point for the curve, with x = 4 and y a unique even number obtained
// from solving curve equation.
func GenPoint() *E521 {
	P := new(E521).getP()
	X := big.NewInt(4)
	num := new(big.Int).Sub(big.NewInt(1), new(big.Int).Exp(X, big.NewInt(2), nil))
	num = num.Mod(num, &P)
	denom := new(big.Int).Add(big.NewInt(1), (new(big.Int).Mul(big.NewInt(376014), new(big.Int).Exp(X, big.NewInt(2), nil))))
	denom = denom.Mod(denom, &P)
	denom = new(big.Int).ModInverse(denom, &P)
	radicand := new(big.Int).Mul(num, denom)

	Y := sqrt(radicand, 0)
	point := E521{
		x: *X,
		y: *Y,
		p: P,
		d: *big.NewInt(-376014),
		r: new(E521).getR(),
		s: "x coord:\n" + X.String() + "\ny coord:\n" + Y.String(),
	}
	return &point

}

// The identity point of the curve (also refered to as "point at infinity").
// Equivalent to 0 in integer group.
func IDPoint() *E521 { return NewE521(*big.NewInt(0), *big.NewInt(1)) }

/*
Gets the opposite value of a point, defined as the following:
if P = (X, Y), opposite of P = (-X, Y).
*/
func (e *E521) getOpposite() *E521 { return NewE521(*e.x.Neg(&e.x), e.y) }

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

	return NewE521(*newX, *newY)
}

/*
EC Multiplication algorithm using the Montgomery Ladder approach to mitigate
power consumption side channel attacks. Mostly constructed around:

	"https://eprint.iacr.org/2014/140.pdf" (pg 4.)

S is a  scalar value to multiply by. S is a private key and should be kept secret.
Returns Curve.E521 point which is result of multiplication.
*/
func (r1 *E521) SecMul(S *big.Int) *E521 {
	r0 := NewE521(*big.NewInt(0), *big.NewInt(1))
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

// func main() {
// 	A := GeneratorPoint()
// 	B := GeneratorPoint().getOpposite()

// 	C := A.Add(B)

// 	// random := generateRandomBigInt()

// 	// pointA = pointA.MultiplyMontgomery(big.NewInt(4))

// 	fmt.Println(A.s)
// 	fmt.Println(B.s)
// 	fmt.Println(C.s)

// }
