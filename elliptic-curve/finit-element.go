package elliptic_curve

import "fmt"

type FieldElement struct {
	order uint64 // field order
	num   uint64 // value of the given element in the field
}

func NewFieldElement(order, num uint64) *FieldElement {
	/*
		init function for FieldElement
	*/
	if num >= order {
		err := fmt.Sprintf("Num not in the range of 0 to %d", order-1)
		panic(err)
	}

	return &FieldElement{
		order: order,
		num:   num,
	}
}

func (f *FieldElement) String() string {
	return fmt.Sprintf("FieldElement{order: %d, num: %d}\n", f.order, f.num)
}

func (f *FieldElement) EqualTo(other *FieldElement) bool {
	return f.order == other.order && f.num == other.num
}
