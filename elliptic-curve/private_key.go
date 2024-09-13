package elliptic_curve

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type PrivateKey struct {
	secret *big.Int
	point  *Point
}

func NewPrivateKey(secret *big.Int) *PrivateKey {
	G := GetGenerator()
	return &PrivateKey{
		secret: secret,
		// public key
		point: G.ScalarMul(secret),
	}
}

func (p *PrivateKey) String() string {
	return fmt.Sprintf("private key hex: {%s}", p.secret)
}

func (p *PrivateKey) GetPublicKey() *Point {
	return p.point
}

func (p *PrivateKey) Sign(z *big.Int) *Signature {
	//(s, r)
	//s = (z + r * e) / k
	// k is a strong random number
	n := GetBitcoinValueN()
	k, err := rand.Int(rand.Reader, n)
	if err != nil {
		panic(fmt.Sprintf("Sign err with rand int: %s", err))
	}
	kField := NewFieldElement(n, k)
	G := GetGenerator()
	// s = (z + r * e) / k
	// r = G * k
	r := G.ScalarMul(k).x.num
	rField := NewFieldElement(n, r)
	eField := NewFieldElement(n, p.secret)
	zField := NewFieldElement(n, z)
	// r*e
	rMulSecret := rField.Multiply(eField)
	// z+r*e
	zAddRMulSecret := zField.Add(rMulSecret)
	// /k
	kInverse := kField.Inverse()
	sField := zAddRMulSecret.Multiply(kInverse)
	/*
	   if s > n / 2 we need to change it to n - s, when doing signature
	   verify, s and n - s are equivalence doing this change is for malleability reasons, detail:
	   https://bitcoin.stackexchange.com/questions/85946/low-s-value-in-bitcoin-signature
	*/
	var opDiv big.Int
	if sField.num.Cmp(opDiv.Div(n, big.NewInt(int64(2)))) > 0 {
		var opSub big.Int
		sField = NewFieldElement(n, opSub.Sub(n, sField.num))
	}

	return &Signature{
		r: NewFieldElement(n, r),
		s: sField,
	}
}
