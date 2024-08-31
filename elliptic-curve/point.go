package elliptic_curve

import (
	"fmt"
	"math/big"
)

type OP_TYPE int

const (
	ADD OP_TYPE = iota
	SUB
	MUL
	DIV
	EXP
)

type Point struct {
	// coefficients of curve
	a *big.Int
	b *big.Int
	// x, y should be the point on the curve
	x *big.Int
	y *big.Int
}

func OpOnBig(x, y *big.Int, opType OP_TYPE) *big.Int {
	var op big.Int
	switch opType {
	case ADD:
		return op.Add(x, y)
	case SUB:
		return op.Sub(x, y)
	case MUL:
		return op.Mul(x, y)
	case DIV:
		return op.Div(x, y)
	case EXP:
		return op.Exp(x, y, nil)
	}

	panic("should not come to here")
}

func NewEllipticCurvePoint(x, y, a, b *big.Int) *Point {
	if x == nil && y == nil {
		return &Point{
			x: x,
			y: y,
			a: a,
			b: b,
		}
	}

	left := OpOnBig(y, big.NewInt(2), EXP)
	x3 := OpOnBig(x, big.NewInt(3), EXP)
	ax := OpOnBig(a, x, MUL)
	right := OpOnBig(OpOnBig(x3, ax, ADD), b, ADD)
	if left.Cmp(right) != 0 {
		err := fmt.Sprintf("Point(%v, %v) is not on the curve with a:%v, b:%v\n", x, y, a, b)
		panic(err)
	}

	return &Point{
		x: x,
		y: y,
		a: a,
		b: b,
	}
}

func (p *Point) Add(other *Point) *Point {
	// check two points are on the same curve
	if p.a.Cmp(other.a) != 0 || p.b.Cmp(other.b) != 0 {
		panic("given two points are not on the same curve")
	}

	if p.x == nil {
		return other
	}

	if other.x == nil {
		return p
	}

	// TODO
	return nil
}

func (p *Point) String() string {
	return fmt.Sprintf("{x: %s, y: %s, a: %s, b: %s}", p.x.String(), p.y.String(), p.a.String(), p.b.String())
}

func (p *Point) Equal(other *Point) bool {
	if p.a.Cmp(other.a) == 0 && p.b.Cmp(other.b) == 0 && p.x.Cmp(other.x) == 0 && p.y.Cmp(other.y) == 0 {
		return true
	}
	return false
}

func (p *Point) NotEqual(other *Point) bool {
	if p.a.Cmp(other.a) != 0 || p.b.Cmp(other.b) != 0 || p.x.Cmp(other.x) != 0 || p.y.Cmp(other.y) != 0 {
		return true
	}
	return false
}
