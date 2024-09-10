package transaction

import (
	"fmt"
	"testing"
)

func TestOpCodeMainTest(t *testing.T) {
	opCode := NewBitcoinOpCode()
	encodeVal := opCode.EncodeNum(-1)
	fmt.Printf("encode -1: %x\n", encodeVal)
	fmt.Printf("decode -1: %d\n", opCode.DecodeNum(encodeVal))

	encodeVal = opCode.EncodeNum(-1234)
	fmt.Printf("encode -1234: %x\n", encodeVal)
	fmt.Printf("decode -1234: %d\n", opCode.DecodeNum(encodeVal))

	encodeVal = opCode.EncodeNum(-32896) // -0x8080
	fmt.Printf("encode -32896: %x\n", encodeVal)
	fmt.Printf("decode -32896: %d\n", opCode.DecodeNum(encodeVal))
}
