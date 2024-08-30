package elliptic_curve

import (
	"fmt"
	"math/big"
)

type FieldElement struct {
	order *big.Int // field order
	num   *big.Int // value of the given element in the field
}

// overflow 64bits integer
// huge number , +, *, ^ => overflow 64bits, we use large or big number

func NewFieldElement(order, num *big.Int) *FieldElement {
	/*
		init function for FieldElement
	*/
	if order.Cmp(num) == -1 {
		err := fmt.Sprintf("Num not in the range of 0 to %d.", order)
		panic(err)
	}

	return &FieldElement{
		order: order,
		num:   num,
	}
}

func (f *FieldElement) String() string {
	return fmt.Sprintf("FieldElement{order: %s, num: %s}\n", f.order.String(), f.num.String())
}

func (f *FieldElement) EqualTo(other *FieldElement) bool {
	return f.order.Cmp(other.order) == 0 && f.num.Cmp(other.num) == 0
}

func (f *FieldElement) checkOrder(other *FieldElement) {
	if f.order.Cmp(other.order) != 0 {
		panic("Add need to do on the field element with the same order.")
	}
}

func (f *FieldElement) Add(other *FieldElement) *FieldElement {
	f.checkOrder(other)

	var op big.Int
	return NewFieldElement(f.order, op.Mod(op.Add(f.num, other.num), f.order))
}

// a, b (a + b) % order = 0, b is called negate of a, b = -a
func (f *FieldElement) Negate() *FieldElement {
	var op big.Int
	return NewFieldElement(f.order, op.Sub(f.order, f.num))
}

func (f *FieldElement) Substract(other *FieldElement) *FieldElement {
	/*
		a, b elemnet of the finite set, c = a - b, given b how can we find c,
		(b + c) % order = a, a - b => (a + (-b)) % order
	*/

	return f.Add(other.Negate())
}

func (f *FieldElement) Multiply(other *FieldElement) *FieldElement {
	f.checkOrder(other)

	// Arithmetic multiplie over modulur of the order
	var op big.Int
	mul := op.Mul(f.num, other.num)
	return NewFieldElement(f.order, op.Mod(mul, f.order))
}

func (f *FieldElement) Power(power *big.Int) *FieldElement {
	// Arithmetic power over modulur of the order
	var op big.Int
	powerRes := op.Exp(f.num, power, nil)
	modRes := op.Mod(powerRes, f.order)
	return NewFieldElement(f.order, modRes)
}

func (f *FieldElement) ScalarMul(val *big.Int) *FieldElement {
	var op big.Int
	res := op.Mul(f.num, val)
	res = op.Mod(res, f.order)
	return NewFieldElement(f.order, res)
}
