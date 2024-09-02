package elliptic_curve

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	/*
		G * k, {G, 2G, ..., nG} n*G->identity generator point
		k is private key
		k*G = Q pub key
		(Q, G) impossible => k

		p = 2^256 - 2^32 - 977
		G(x,y)
		Gx = 0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798
		Gy = 0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8
		n = 0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141
		y^2 = x^3 + 7
	*/
	var op big.Int
	twoExp256 := op.Exp(big.NewInt(2), big.NewInt(256), nil)
	var op1 big.Int
	twoExp32 := op1.Exp(big.NewInt(2), big.NewInt(32), nil)
	var op2 big.Int
	p := op2.Sub(twoExp256, twoExp32)
	var op3 big.Int
	pp := op3.Sub(p, big.NewInt(977))
	fmt.Printf("pp is %s\n", pp)
	Gx := new(big.Int)
	Gx.SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
	Gy := new(big.Int)
	Gy.SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)

	x1 := NewFieldElement(pp, Gx)
	y1 := NewFieldElement(pp, Gy)
	a := NewFieldElement(pp, big.NewInt(0))
	b := NewFieldElement(pp, big.NewInt(7))
	G := NewEllipticCurvePoint(x1, y1, a, b)
	fmt.Printf("G is on elliptic curve with value is %s\n", G)

	G = S256Point(Gx, Gy)
	n := new(big.Int)
	n.SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	fmt.Printf("n*G is :%s\n", G.ScalarMul(n))
}

func TestCheckPointOnCurve(t *testing.T) {
	var (
		x *FieldElement
		y *FieldElement
		a = NewFieldElement(big.NewInt(223), big.NewInt(0))
		b = NewFieldElement(big.NewInt(223), big.NewInt(7))
	)

	x = NewFieldElement(big.NewInt(223), big.NewInt(192))
	y = NewFieldElement(big.NewInt(223), big.NewInt(105))
	assert.NotPanics(t, func() {
		NewEllipticCurvePoint(x, y, a, b)
	})

	y = NewFieldElement(big.NewInt(223), big.NewInt(106))
	assert.Panics(t, func() {
		NewEllipticCurvePoint(x, y, a, b)
	})
}

func TestAddIdentity(t *testing.T) {
	var (
		x *FieldElement
		y *FieldElement
		a = NewFieldElement(big.NewInt(223), big.NewInt(0))
		b = NewFieldElement(big.NewInt(223), big.NewInt(7))
	)
	x = NewFieldElement(big.NewInt(223), big.NewInt(192))
	y = NewFieldElement(big.NewInt(223), big.NewInt(105))
	p := NewEllipticCurvePoint(x, y, a, b)
	fmt.Printf("p is %s\n", p)

	identity := NewEllipticCurvePoint(nil, nil, a, b)
	assert.True(t, p.Add(identity).Equal(p))
}

func TestAddVertical(t *testing.T) {
	var (
		x *FieldElement
		y *FieldElement
		a = NewFieldElement(big.NewInt(223), big.NewInt(0))
		b = NewFieldElement(big.NewInt(223), big.NewInt(7))
	)
	x = NewFieldElement(big.NewInt(223), big.NewInt(192))
	y = NewFieldElement(big.NewInt(223), big.NewInt(105))
	p1 := NewEllipticCurvePoint(x, y, a, b)
	yNeg := y.Negate()
	p2 := NewEllipticCurvePoint(x, yNeg, a, b)
	fmt.Printf("addition of points on vertial line over finite field is %s\n", p1.Add(p2))
}

// func TestAddSelf(t *testing.T) {
// 	// C = A + A
// 	A := NewEllipticCurvePoint(big.NewInt(-1), big.NewInt(-1), big.NewInt(5), big.NewInt(7))
// 	C := NewEllipticCurvePoint(big.NewInt(18), big.NewInt(77), big.NewInt(5), big.NewInt(7))
// 	assert.True(t, A.Add(A).Equal(C))
// }

func TestAdd(t *testing.T) {
	var (
		x *FieldElement
		y *FieldElement
		a = NewFieldElement(big.NewInt(223), big.NewInt(0))
		b = NewFieldElement(big.NewInt(223), big.NewInt(7))
	)
	x = NewFieldElement(big.NewInt(223), big.NewInt(192))
	y = NewFieldElement(big.NewInt(223), big.NewInt(105))
	p1 := NewEllipticCurvePoint(x, y, a, b)
	x = NewFieldElement(big.NewInt(223), big.NewInt(17))
	y = NewFieldElement(big.NewInt(223), big.NewInt(56))
	p2 := NewEllipticCurvePoint(x, y, a, b)
	p3 := NewEllipticCurvePoint(NewFieldElement(big.NewInt(223), big.NewInt(170)), NewFieldElement(big.NewInt(223), big.NewInt(142)), a, b)
	assert.True(t, p1.Add(p2).Equal(p3))
}

func TestScalarMul(t *testing.T) {
	// 2 * (192, 105)
	p1 := newPoint(t, 192, 105, 0, 7)
	res := newPoint(t, 49, 71, 0, 7)
	assert.True(t, p1.ScalarMul(big.NewInt(2)).Equal(res))

	// 2 * (47, 71)
	p2 := newPoint(t, 47, 71, 0, 7)
	res = newPoint(t, 36, 111, 0, 7)
	assert.True(t, p2.ScalarMul(big.NewInt(2)).Equal(res))
}

func newPoint(t *testing.T, x, y, a, b int64) *Point {
	xF := NewFieldElement(big.NewInt(223), big.NewInt(x))
	yF := NewFieldElement(big.NewInt(223), big.NewInt(y))
	aF := NewFieldElement(big.NewInt(223), big.NewInt(a))
	bF := NewFieldElement(big.NewInt(223), big.NewInt(b))
	var p *Point
	assert.NotPanics(t, func() {
		p = NewEllipticCurvePoint(xF, yF, aF, bF)
	})
	return p
}
