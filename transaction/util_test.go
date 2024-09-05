package transaction

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConversionBetweenBigEndianAndLittleEndian(t *testing.T) {
	p := new(big.Int)
	p.SetString("12345678", 16)
	bytes := p.Bytes()
	fmt.Printf("bytes for 0x12345678 is %x\n", bytes)

	littleEndianByte := BigIntToLittleEndian(p, LITTLE_ENDIAN_4_BYTES)
	fmt.Printf("little endian for 0x12345678 is %x\n", littleEndianByte)

	littleEndianByteToInt64 := LittleEndianToBigInt(littleEndianByte, LITTLE_ENDIAN_4_BYTES)
	fmt.Printf("little endian bytes into int is %x\n", littleEndianByteToInt64)

	assert.Equal(t, p.Cmp(littleEndianByteToInt64), 0)
}
