package elliptic_curve

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSEC(t *testing.T) {
	// 0xdeadbeef54321
	p := new(big.Int)
	p.SetString("deadbeef54321", 16)
	privateKey := NewPrivateKey(p)
	pubKey := privateKey.GetPublicKey()
	fmt.Printf("public key is %s\n", pubKey)

	secBinUnCompressed := new(big.Int)
	secBinUnCompressed.SetString(pubKey.Sec(false), 16)
	unUnCompressedDecode := ParseSEC(secBinUnCompressed.Bytes())
	fmt.Printf("decode sec uncompressed format: %s\n", unUnCompressedDecode)
	assert.True(t, pubKey.Equal(unUnCompressedDecode))

	secBinCompressed := new(big.Int)
	secBinCompressed.SetString(pubKey.Sec(true), 16)
	compressedDecode := ParseSEC(secBinCompressed.Bytes())
	fmt.Printf("decode sec compressed format: %s\n", compressedDecode)
	assert.True(t, pubKey.Equal(compressedDecode))
}
