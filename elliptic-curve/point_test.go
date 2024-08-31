package elliptic_curve

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
