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

func S256Field(num *big.Int) *FieldElement {
	// 2^256 - 2^32 - 977
	var op big.Int
	twoExp256 := op.Exp(big.NewInt(2), big.NewInt(256), nil)
	var op1 big.Int
	twoExp32 := op1.Exp(big.NewInt(2), big.NewInt(32), nil)
	var op2 big.Int
	p := op2.Sub(twoExp256, twoExp32)
	var op3 big.Int
	pp := op3.Sub(p, big.NewInt(977))
	return NewFieldElement(pp, num)
}

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
	return fmt.Sprintf("FieldElement{order: %x, num: %x}\n", f.order, f.num)
}

func (f *FieldElement) EqualTo(other *FieldElement) bool {
	return f.order.Cmp(other.order) == 0 && f.num.Cmp(other.num) == 0
}

func (f *FieldElement) checkOrder(other *FieldElement) {
	if f.order.Cmp(other.order) != 0 {
		panic("need to do on the field element with the same order.")
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
	// k ^ (p - 1) % p = 1, power > p - 1 => power % (p - 1)
	var op big.Int
	t := op.Mod(power, op.Sub(f.order, big.NewInt(1)))
	powerRes := op.Exp(f.num, t, f.order)
	// modRes := op.Mod(powerRes, f.order)
	return NewFieldElement(f.order, powerRes)
}

func (f *FieldElement) ScalarMul(val *big.Int) *FieldElement {
	var op big.Int
	res := op.Mul(f.num, val)
	res = op.Mod(res, f.order)
	return NewFieldElement(f.order, res)
}

func (f *FieldElement) Divide(other *FieldElement) *FieldElement {
	f.checkOrder(other)
	// a / b => a . b ^ (p -2)
	var op big.Int
	otherReverse := other.Power(op.Sub(f.order, big.NewInt(int64(2))))
	return f.Multiply(otherReverse)
}

func (f *FieldElement) Inverse() *FieldElement {
	var op big.Int
	return f.Power(op.Sub(f.order, big.NewInt(2)))
}

/*
04x
y^2 = x^3 + 7 => y
w^2 = v => w
1. p % 4 == 3, (p+1)/4 is integer
2. Fermat's little w^(p-1)%p == 1, w^2 * 1 == w^2 * w^(p-1) => w^2=w^(p-1)
3. p is prime => p is odd, (p+1)/2 is integer
4. w^(p+1) = w^((p+1)/2) * w^((p+1)/2) => [w^((p+1)/2)]^2 = w^(p+1) = w^2
w^((p+1)/2) == w
5. w^(2*(p+1)/4) = w^((p+1)/2) => (w^2)^((p+1)/4) = w^((p+1)/2) = w
6. (w^2)^((p+1)/4) = w^((p+1)/2) = w = v^((p+1)/4) == w

v^((p+1)/4) == w
*/
func (f *FieldElement) Sqrt() *FieldElement {
	// make sure (p+1) % 4 == 0
	var opAdd big.Int
	orderAddOne := opAdd.Add(f.order, big.NewInt(1))
	var opMod big.Int
	modRes := opMod.Mod(orderAddOne, big.NewInt(4))
	if modRes.Cmp(big.NewInt(0)) != 0 {
		panic("order plus one mod 4 is not 0")
	}

	var opDiv big.Int
	return f.Power(opDiv.Div(orderAddOne, big.NewInt(4)))
}
