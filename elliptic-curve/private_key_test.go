package elliptic_curve

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	e := big.NewInt(12345)
	z := new(big.Int)
	z.SetBytes(Hash256("Testing my Signing"))

	privateKey := NewPrivateKey(e)
	sig := privateKey.Sign(z)
	fmt.Printf("sig is %s\n", sig)

	pubKey := privateKey.GetPublicKey()
	n := GetBitcoinValueN()
	zField := NewFieldElement(n, z)
	res := pubKey.Verify(zField, sig)
	assert.True(t, res)
}
