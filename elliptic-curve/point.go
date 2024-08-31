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
	a *FieldElement
	b *FieldElement
	// x, y should be the point on the curve
	x *FieldElement
	y *FieldElement
}

func OpOnBig(x, y *FieldElement, scalar *big.Int, opType OP_TYPE) *FieldElement {
	switch opType {
	case ADD:
		return x.Add(y)
	case SUB:
		return x.Substract(y)
	case MUL:
		if y != nil {
			return x.Multiply(y)
		}
		if scalar != nil {
			return x.ScalarMul(scalar)
		}
		panic("error in multiply")
	case DIV:
		return x.Divide(y)
	case EXP:
		if scalar == nil {
			panic("scalar should not be nil for EXP")
		}
		return x.Power(scalar)
	}

	panic("should not come to here")
}

func NewEllipticCurvePoint(x, y, a, b *FieldElement) *Point {
	if x == nil && y == nil {
		return &Point{
			x: x,
			y: y,
			a: a,
			b: b,
		}
	}

	left := OpOnBig(y, nil, big.NewInt(2), EXP)
	x3 := OpOnBig(x, nil, big.NewInt(3), EXP)
	ax := OpOnBig(a, x, nil, MUL)
	right := OpOnBig(OpOnBig(x3, ax, nil, ADD), b, nil, ADD)
	if !left.EqualTo(right) {
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

/*
G != identity {G, 2*G, ..., n*G} n*G identity
k * G => Q easy  G, Q => k impossible
k*G => G + G + ... + G
k => 13(1101) (2^3 + 2^2 + 2^0) * G => 2^3*G + 2^2*G + 2^0*G
=> (G<<3) + (G<<2) + (G<<0) k has t 1s in binary form, we can do t times of addition
1 trillion, 40 bits in binary form
we at most do 40 times of addition => 1 trillion times
*/
func (p *Point) ScalarMul(scalar *big.Int) *Point {
	if scalar == nil {
		panic("scalar can not be nil")
	}

	// 13 => "1101"
	binaryFrom := fmt.Sprintf("%b", scalar)
	current := p
	result := NewEllipticCurvePoint(nil, nil, p.a, p.b)
	for i := len(binaryFrom) - 1; i >= 0; i-- {
		if binaryFrom[i] == '1' {
			result = result.Add(current)
		}
		// left shift by 1 place, just like add to self
		current = current.Add(current)
	}
	return result
}

func (p *Point) Add(other *Point) *Point {
	// check two points are on the same curve
	if !p.a.EqualTo(other.a) || !p.b.EqualTo(other.b) {
		panic("given two points are not on the same curve")
	}

	if p.x == nil {
		return other
	}

	if other.x == nil {
		return p
	}

	// points are on the verical A(x,y), b(x,-y)
	zero := NewFieldElement(p.x.order, big.NewInt(0))
	if p.x.EqualTo(other.x) && OpOnBig(p.y, other.y, nil, ADD).EqualTo(zero) {
		return &Point{
			x: nil,
			y: nil,
			a: p.a,
			b: p.b,
		}
	}

	// find slope of line AB
	// x1 -> p.x, y1 -> p.y, x2 -> other.x, y2 -> other.y
	var numerator *FieldElement
	var denominator *FieldElement
	if p.x.EqualTo(other.x) && p.y.EqualTo(other.y) {
		// slope = (3*x^2+a) / 2y
		xSqrt := OpOnBig(p.x, nil, big.NewInt(2), EXP)
		threeXSqrt := OpOnBig(xSqrt, nil, big.NewInt(3), MUL)
		numerator = OpOnBig(threeXSqrt, p.a, nil, ADD)
		denominator = OpOnBig(p.y, nil, big.NewInt(2), MUL)
	} else {
		// s= (y2-y1) / (x2-x1)
		numerator = OpOnBig(other.y, p.y, nil, SUB)
		denominator = OpOnBig(other.x, p.x, nil, SUB)
	}

	slope := OpOnBig(numerator, denominator, nil, DIV)
	// s^2
	slopeSqrt := OpOnBig(slope, nil, big.NewInt(2), EXP)
	// x3 = s^2 - x1 - x2
	x3 := OpOnBig(OpOnBig(slopeSqrt, p.x, nil, SUB), other.x, nil, SUB)
	// x3 - x1
	x3Minusx1 := OpOnBig(x3, p.x, nil, SUB)
	// y3 = s(x3 - x1) + y1
	y3 := OpOnBig(OpOnBig(slope, x3Minusx1, nil, MUL), p.y, nil, ADD)
	// -y3
	minusy3 := OpOnBig(y3, nil, big.NewInt(-1), MUL)

	return &Point{
		x: x3,
		y: minusy3,
		a: p.a,
		b: p.b,
	}
}

func (p *Point) String() string {
	xString := "nil"
	yString := "nil"
	if p.x != nil {
		xString = p.x.String()
	}
	if p.y != nil {
		yString = p.y.String()
	}
	return fmt.Sprintf("{x: %s, y: %s, a: %s, b: %s}", xString, yString, p.a.String(), p.b.String())
}

func (p *Point) Equal(other *Point) bool {
	if p.a.EqualTo(other.a) && p.b.EqualTo(other.b) && p.x.EqualTo(other.x) && p.y.EqualTo(other.y) {
		return true
	}
	return false
}

func (p *Point) NotEqual(other *Point) bool {
	if !p.a.EqualTo(other.a) || !p.b.EqualTo(other.b) || !p.x.EqualTo(other.x) || !p.y.EqualTo(other.y) {
		return true
	}
	return false
}
