package elliptic_curve

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SolvField19MultiplieSet() {

}

func TestSolvField19MultiplieSet(t *testing.T) {
	// randomly select a num from 1 to 18
	min := 1
	max := 18
	k := rand.Intn(max-min) + 1
	fmt.Printf("randomly select k is %d\n", k)
	element := NewFieldElement(big.NewInt(19), big.NewInt(int64(k)))
	for i := 0; i < 19; i++ {
		fmt.Printf("element %d multiplie with %d is %v\n", k, i, element.ScalarMul(big.NewInt(int64(i))))
	}
}

/*
p is field order, p = 7, 11, 17, 19, 31
select any element k in the field with order p, compute k ^ (p-1) what is the result
=> {1 ^ (p-1) % p, 2 ^ (p-1) % p, ......, (p-1) ^ (p-1) % p}

for any element k in the field with order => k ^ (p-1) % p == 1
*/
func TestComputeFieldOrderPower(t *testing.T) {
	orders := []int{7, 11, 17, 19, 31}
	for _, p := range orders {
		fmt.Printf("value of p is %d\n", p)
		for i := 1; i < p; i++ {
			elem := NewFieldElement(big.NewInt(int64(p)), big.NewInt(int64(i)))
			fmt.Printf("for element %v, its power of p - 1 is %v\n", elem, elem.Power(big.NewInt(int64(p-1))))
		}
	}
}

func TestDivideFiniteElement(t *testing.T) {
	f2 := NewFieldElement(big.NewInt(19), big.NewInt(2))
	f7 := NewFieldElement(big.NewInt(19), big.NewInt(7))
	fmt.Printf("field element 2 / 7 is %v\n", f2.Divide(f7))

	f46 := NewFieldElement(big.NewInt(57), big.NewInt(46))
	fmt.Printf("field element 46 * 46 with order 57 is %v\n", f46.Multiply(f46))
	fmt.Printf("field element 46 with power of 58 is %v\n", f46.Power(big.NewInt(58)))
	// 58 % 56 = 2
}

func TestSqrtFiniteElement(t *testing.T) {
	x := new(big.Int)
	x.SetString("7211a824f55b505228e4c3d5194c1fcfaa15a456abdf37f9b9d97a4040afc073", 16)
	// y^2 = x^3 + 7
	y2 := S256Field(x).Power(big.NewInt(3)).Add(S256Field(big.NewInt(7)))
	y := y2.Sqrt()
	fmt.Printf("y value of given x is %s\n", y)

	// check (x,y) is on the curve
	assert.NotPanics(t, func() {
		NewEllipticCurvePoint(S256Field(x), y, S256Field(big.NewInt(0)), S256Field(big.NewInt(7)))
	})
}
