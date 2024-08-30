package main

import (
	ecc "elliptic_curve"
	"fmt"
	"math/rand"
)

func SolvField19MultiplieSet() {
	// randomly select a num from 1 to 18
	min := 1
	max := 18
	k := rand.Intn(max-min) + 1
	fmt.Printf("randomly select k is %d\n", k)
	element := ecc.NewFieldElement(19, uint64(k))
	for i := 0; i < 19; i++ {
		fmt.Printf("element %d multiplie with %d is %v\n", k, i, element.ScalarMul(uint64(i)))
	}
}

func main() {
	SolvField19MultiplieSet()
}
