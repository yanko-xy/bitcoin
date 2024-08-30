package main

import (
	ecc "elliptic_curve"
	"fmt"
)

func main() {
	f44 := ecc.NewFieldElement(57, 44)
	f33 := ecc.NewFieldElement(57, 33)
	res := f44.Add(f33)
	fmt.Printf("Field element 44 add to field field element 33 is %v\n", res)
	// -44 negate of 44 is 57 - 44 is 13
	fmt.Printf("Negate of field element 44 is %v\n", res.Negate())
	fmt.Printf("Field element 44 - 33 is %v\n", f44.Substract(f33))
	fmt.Printf("Field element 33 - 44 is %v\n", f33.Substract(f44))

	// check (46 + 44) % order == 33
	fmt.Printf("check 46 + 44 over modulur of 57 is %d\n", (46+44)%57)
	// f33 - f44 = f46 => f46 + f44 = f33
	f46 := ecc.NewFieldElement(57, 46)
	fmt.Printf("Field element 46 + 44 is %v\n", f46.Add(f44))
}
