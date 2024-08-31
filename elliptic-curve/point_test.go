package elliptic_curve

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckPointOnCurve(t *testing.T) {
	/*
		check point(-1, -1) on curve of y^2 = x^3 + 5x + 7
	*/
	assert.NotPanics(t, func() {
		NewEllipticCurvePoint(big.NewInt(-1), big.NewInt(-1), big.NewInt(5), big.NewInt(7))
	})
	/*
		check point(-1, -2) on curve of y^2 = x^3 + 5x + 7
	*/
	assert.Panics(t, func() {
		NewEllipticCurvePoint(big.NewInt(-1), big.NewInt(-2), big.NewInt(5), big.NewInt(7))
	})

	/*
		(2, 4) (18, 77) (5, 7) on the curve of y^2 = x^3 + 5x + 7
	*/
	assert.Panics(t, func() {
		NewEllipticCurvePoint(big.NewInt(2), big.NewInt(4), big.NewInt(5), big.NewInt(7))
	})
	assert.NotPanics(t, func() {
		NewEllipticCurvePoint(big.NewInt(18), big.NewInt(77), big.NewInt(5), big.NewInt(7))
	})
	assert.Panics(t, func() {
		NewEllipticCurvePoint(big.NewInt(5), big.NewInt(7), big.NewInt(5), big.NewInt(7))
	})
}

func TestAddIdentity(t *testing.T) {
	p := NewEllipticCurvePoint(big.NewInt(-1), big.NewInt(-1), big.NewInt(5), big.NewInt(7))
	fmt.Printf("p is %s\n", p)

	identity := NewEllipticCurvePoint(nil, nil, big.NewInt(5), big.NewInt(7))
	assert.True(t, p.Add(identity).Equal(p))
}

func TestAddVertical(t *testing.T) {
	// (-1, -1) (-1, 1)
	p1 := NewEllipticCurvePoint(big.NewInt(-1), big.NewInt(-1), big.NewInt(5), big.NewInt(7))
	p2 := NewEllipticCurvePoint(big.NewInt(-1), big.NewInt(1), big.NewInt(5), big.NewInt(7))
	indentity := NewEllipticCurvePoint(nil, nil, big.NewInt(5), big.NewInt(7))
	assert.True(t, p1.Add(p2).Equal(indentity))
}

func TestAddSelf(t *testing.T) {
	// C = A + A
	A := NewEllipticCurvePoint(big.NewInt(-1), big.NewInt(-1), big.NewInt(5), big.NewInt(7))
	C := NewEllipticCurvePoint(big.NewInt(18), big.NewInt(77), big.NewInt(5), big.NewInt(7))
	assert.True(t, A.Add(A).Equal(C))
}

func TestAdd(t *testing.T) {
	// C = A(2, 5) + B(-1, -1)
	A := NewEllipticCurvePoint(big.NewInt(2), big.NewInt(5), big.NewInt(5), big.NewInt(7))
	B := NewEllipticCurvePoint(big.NewInt(-1), big.NewInt(-1), big.NewInt(5), big.NewInt(7))
	C := NewEllipticCurvePoint(big.NewInt(3), big.NewInt(-7), big.NewInt(5), big.NewInt(7))
	assert.True(t, A.Add(B).Equal(C))
}
