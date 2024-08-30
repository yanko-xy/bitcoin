package main

import (
	ecc "elliptic_curve"
	"fmt"
	"math/big"
	"math/rand"
)

func SolvField19MultiplieSet() {
	// randomly select a num from 1 to 18
	min := 1
	max := 18
	k := rand.Intn(max-min) + 1
	fmt.Printf("randomly select k is %d\n", k)
	element := ecc.NewFieldElement(big.NewInt(19), big.NewInt(int64(k)))
	for i := 0; i < 19; i++ {
		fmt.Printf("element %d multiplie with %d is %v\n", k, i, element.ScalarMul(big.NewInt(int64(i))))
	}
}

func main() {
	f44 := ecc.NewFieldElement(big.NewInt(57), big.NewInt(44))
	f33 := ecc.NewFieldElement(big.NewInt(57), big.NewInt(33))
	res := f44.Add(f33)
	fmt.Printf("field element 44 add to field element 33 is %v\n", res)
	// -44 => 57 - 44 = 13
	fmt.Printf("negate to field element 44 is %v\n", f44.Negate())

	// 44 - 33
	fmt.Printf("field element  44 - 33 is %v\n", f44.Substract(f33))
	// 33 - 44
	fmt.Printf("field element 33 - 44 is %v\n", f33.Substract(f44))
	// check (46 + 44) % 57 == 33
	fmt.Printf("check 46 + 43 over modulur 57 %v\n", (46+44)%57)
	f46 := ecc.NewFieldElement(big.NewInt(57), big.NewInt(46))
	fmt.Printf("field element 46 + 44 is %v\n", f46.Add(f44))

	SolvField19MultiplieSet()
}
