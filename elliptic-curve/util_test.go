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

func TestEncodeBase58(t *testing.T) {
	val := new(big.Int)
	val.SetString("7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d", 16)
	fmt.Printf("base58 encoding is %s\n", EncodeBase58(val.Bytes()))

	val.SetString("eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c", 16)
	fmt.Printf("base58 encoding is %s\n", EncodeBase58(val.Bytes()))

	val.SetString("c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6", 16)
	fmt.Printf("base58 encoding is %s\n", EncodeBase58(val.Bytes()))
}
