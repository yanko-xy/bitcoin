package elliptic_curve

import (
	"fmt"
	"math/big"
	"testing"
)

func TestSIgnatureDer(t *testing.T) {
	r := new(big.Int)
	r.SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
	rField := S256Field(r)
	s := new(big.Int)
	s.SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)
	sField := S256Field(s)
	sig := NewSignature(rField, sField)
	derEncode := sig.Der()
	fmt.Printf("der encoding for signature is %x\n", derEncode)

	sig2 := ParseSigBin(derEncode)
	fmt.Printf("signature parsed from raw binary data is %s\n", sig2)
}
