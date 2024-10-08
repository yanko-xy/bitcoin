package elliptic_curve

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	/*
		G * k, {G, 2G, ..., nG} n*G->identity generator point
		k is private key
		k*G = Q pub key
		(Q, G) impossible => k

		p = 2^256 - 2^32 - 977
		G(x,y)
		Gx = 0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798
		Gy = 0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8
		n = 0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141
		y^2 = x^3 + 7

		ECDSA -> Elliptic Curve Digital Signature Algorithm
		1. select a scalar e, compute p = e*G, P is public key you can release out
		2. e is private key you need keep, randomly select two field number u, v
		compute k = u + e*v, k need to keep secret
		3. compute R = k*G = (u + e*v) * G = u*G + v*(e*G) = u*G + v*P, take the
		x coordinate out, take this value as r
		4. owner of e, generate a text message in format of string, hash it 256 bits integer(sha256, md5),
		called the hash result as z
		5. compute number of s = (z + r*e)/k (base on modulur p)
		6. release (z, s, r) as signature of the key owner
		7. any one who want to verify message z is created by owner of e:
			1). compute u = z/s, v=r/s
			2). compute u*G + v*P = (z/s)*G + (r/s)*P = (z/s)*G + (r/s)*eG
			= (z/s)*G + (r*e/s)*G = ((z+r*e)/s)*G = k*G= R'
			3). take the x coordinate of R', compare with r
			if the same => verify the message z is created by owner of e

			(z, s, r, P) / is multiply inverse, it is not the normal arithmetic op
	*/
	var op big.Int
	twoExp256 := op.Exp(big.NewInt(2), big.NewInt(256), nil)
	var op1 big.Int
	twoExp32 := op1.Exp(big.NewInt(2), big.NewInt(32), nil)
	var op2 big.Int
	p := op2.Sub(twoExp256, twoExp32)
	var op3 big.Int
	pp := op3.Sub(p, big.NewInt(977))
	fmt.Printf("pp is %s\n", pp)
	Gx := new(big.Int)
	Gx.SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
	Gy := new(big.Int)
	Gy.SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)

	x1 := NewFieldElement(pp, Gx)
	y1 := NewFieldElement(pp, Gy)
	a := NewFieldElement(pp, big.NewInt(0))
	b := NewFieldElement(pp, big.NewInt(7))
	G := NewEllipticCurvePoint(x1, y1, a, b)
	fmt.Printf("G is on elliptic curve with value is %s\n", G)

	G = S256Point(Gx, Gy)
	n := new(big.Int)
	n.SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	fmt.Printf("n*G is :%s\n", G.ScalarMul(n))
}

func TestCheckPointOnCurve(t *testing.T) {
	var (
		x *FieldElement
		y *FieldElement
		a = NewFieldElement(big.NewInt(223), big.NewInt(0))
		b = NewFieldElement(big.NewInt(223), big.NewInt(7))
	)

	x = NewFieldElement(big.NewInt(223), big.NewInt(192))
	y = NewFieldElement(big.NewInt(223), big.NewInt(105))
	assert.NotPanics(t, func() {
		NewEllipticCurvePoint(x, y, a, b)
	})

	y = NewFieldElement(big.NewInt(223), big.NewInt(106))
	assert.Panics(t, func() {
		NewEllipticCurvePoint(x, y, a, b)
	})
}

func TestAddIdentity(t *testing.T) {
	var (
		x *FieldElement
		y *FieldElement
		a = NewFieldElement(big.NewInt(223), big.NewInt(0))
		b = NewFieldElement(big.NewInt(223), big.NewInt(7))
	)
	x = NewFieldElement(big.NewInt(223), big.NewInt(192))
	y = NewFieldElement(big.NewInt(223), big.NewInt(105))
	p := NewEllipticCurvePoint(x, y, a, b)
	fmt.Printf("p is %s\n", p)

	identity := NewEllipticCurvePoint(nil, nil, a, b)
	assert.True(t, p.Add(identity).Equal(p))
}

func TestAddVertical(t *testing.T) {
	var (
		x *FieldElement
		y *FieldElement
		a = NewFieldElement(big.NewInt(223), big.NewInt(0))
		b = NewFieldElement(big.NewInt(223), big.NewInt(7))
	)
	x = NewFieldElement(big.NewInt(223), big.NewInt(192))
	y = NewFieldElement(big.NewInt(223), big.NewInt(105))
	p1 := NewEllipticCurvePoint(x, y, a, b)
	yNeg := y.Negate()
	p2 := NewEllipticCurvePoint(x, yNeg, a, b)
	fmt.Printf("addition of points on vertial line over finite field is %s\n", p1.Add(p2))
}

// func TestAddSelf(t *testing.T) {
// 	// C = A + A
// 	A := NewEllipticCurvePoint(big.NewInt(-1), big.NewInt(-1), big.NewInt(5), big.NewInt(7))
// 	C := NewEllipticCurvePoint(big.NewInt(18), big.NewInt(77), big.NewInt(5), big.NewInt(7))
// 	assert.True(t, A.Add(A).Equal(C))
// }

func TestAdd(t *testing.T) {
	var (
		x *FieldElement
		y *FieldElement
		a = NewFieldElement(big.NewInt(223), big.NewInt(0))
		b = NewFieldElement(big.NewInt(223), big.NewInt(7))
	)
	x = NewFieldElement(big.NewInt(223), big.NewInt(192))
	y = NewFieldElement(big.NewInt(223), big.NewInt(105))
	p1 := NewEllipticCurvePoint(x, y, a, b)
	x = NewFieldElement(big.NewInt(223), big.NewInt(17))
	y = NewFieldElement(big.NewInt(223), big.NewInt(56))
	p2 := NewEllipticCurvePoint(x, y, a, b)
	p3 := NewEllipticCurvePoint(NewFieldElement(big.NewInt(223), big.NewInt(170)), NewFieldElement(big.NewInt(223), big.NewInt(142)), a, b)
	assert.True(t, p1.Add(p2).Equal(p3))
}

func TestScalarMul(t *testing.T) {
	// 2 * (192, 105)
	p1 := newPoint(t, 192, 105, 0, 7)
	res := newPoint(t, 49, 71, 0, 7)
	assert.True(t, p1.ScalarMul(big.NewInt(2)).Equal(res))

	// 2 * (47, 71)
	p2 := newPoint(t, 47, 71, 0, 7)
	res = newPoint(t, 36, 111, 0, 7)
	assert.True(t, p2.ScalarMul(big.NewInt(2)).Equal(res))
}

func TestVerify(t *testing.T) {
	/*
		        P(x,y)= (
		             0x4519fac3d910ca7e7138f7013706f619fa8f033e6ec6e09370ea38cee6a7574，
		             82b51eab8c27c66e26c858a079bcdf4f1ada34cec420cafc7eac1a42216fb6c4
				)

				z: 0xbc62d4b80d9e36da29c16c5d4d9f11731f36052c72401a76c23c0fb5a9b74423
				r: 37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6
				s: 8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec

				R is public key used to verify the message, z is hashed message,
				s is generated by using z and private key e,
				Verify should return true
	*/
	n := new(big.Int)
	n.SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)

	zVal := new(big.Int)
	zVal.SetString("bc62d4b80d9e36da29c16c5d4d9f11731f36052c72401a76c23c0fb5a9b74423", 16)
	zField := NewFieldElement(n, zVal)

	rVal := new(big.Int)
	rVal.SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
	rField := NewFieldElement(n, rVal)

	sVal := new(big.Int)
	sVal.SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)
	sField := NewFieldElement(n, sVal)

	//public key
	px := new(big.Int)
	px.SetString("4519fac3d910ca7e7138f7013706f619fa8f033e6ec6e09370ea38cee6a7574", 16)
	py := new(big.Int)
	py.SetString("82b51eab8c27c66e26c858a079bcdf4f1ada34cec420cafc7eac1a42216fb6c4", 16)
	point := S256Point(px, py)

	sig := NewSignature(rField, sField)
	verifyRes := point.Verify(zField, sig)
	assert.True(t, verifyRes)
}

/*
When a user of bitcoin create a wallet, he/she need to publish the address of his or her wallet and
he can receive or sending funds to others. Wallet address is actually public key we created at last section and
encode in some kind of format called SEC(standards for efficient cryptography).
SEC format has two forms they are uncompressed and compressed,let's check the uncompressed format first.
For a public key P=(x,y), the coordinate of x and y are 32 bytes integer, we use the following steps to encode the key into uncompressed SEC format:

1, the beginning byte set to 0x04,
2, turn the x value into big-endian and append it after byte 0x04
3, turn the y value into big-endian and append it after the end of x

Let's see an example of SEC uncompressed format:

047211a824f55b505228e4c3d5194c1fcfaa15a456abdf37f9b9d97a4040afc073dee6c89064984f03385237d92167c13e236446b417ab79a0fcae412ae3316b77

let's split the chunk of data into three parts:

1. beginning byte 0x04
2. x value in big-endian: 7211a824f55b505228e4c3d5194c1fcfaa15a456abdf37f9b9d97a4040afc073
3. y value in big-endian: dee6c89064984f03385237d92167c13e236446b417ab79a0fcae412ae3316b77
let's check whether the point with the given x,y is on the bitcoin curve:
*/
func TestUmcompressedSec(t *testing.T) {
	privateKey := NewPrivateKey(big.NewInt(5000))
	publicKey := privateKey.GetPublicKey()
	_, secBytes := publicKey.Sec(false)
	fmt.Printf("sec uncompress format for 5000*G is %s\n", secBytes)

	// 2018 ^ 5
	var expOp big.Int
	privateKey = NewPrivateKey(expOp.Exp(big.NewInt(2018), big.NewInt(5), nil))
	publicKey = privateKey.GetPublicKey()
	_, secBytes = publicKey.Sec(false)
	fmt.Printf("sec uncompress format for 2018^5*G is %s\n", secBytes)

	// 0xdeadbeef123456
	p := new(big.Int)
	p.SetString("deadbeef123456", 16)
	privateKey = NewPrivateKey(p)
	publicKey = privateKey.GetPublicKey()
	_, secBytes = publicKey.Sec(false)
	fmt.Printf("sec uncompress format for 50xdeadbeef123456*G is %s\n", secBytes)
}

func TestCompressedSec(t *testing.T) {
	privateKey := NewPrivateKey(big.NewInt(5001))
	publicKey := privateKey.GetPublicKey()
	_, secBytes := publicKey.Sec(true)
	fmt.Printf("sec compress format for 5000*G is %s\n", secBytes)

	// 2018 ^ 5
	var expOp big.Int
	privateKey = NewPrivateKey(expOp.Exp(big.NewInt(2018), big.NewInt(5), nil))
	publicKey = privateKey.GetPublicKey()
	_, secBytes = publicKey.Sec(true)
	fmt.Printf("sec compress format for 2018^5*G is %s\n", secBytes)

	// 0xdeadbeef123456
	p := new(big.Int)
	p.SetString("deadbeef123456", 16)
	privateKey = NewPrivateKey(p)
	publicKey = privateKey.GetPublicKey()
	_, secBytes = publicKey.Sec(true)
	fmt.Printf("sec compress format for 50xdeadbeef123456*G is %s\n", secBytes)
}

func newPoint(t *testing.T, x, y, a, b int64) *Point {
	xF := NewFieldElement(big.NewInt(223), big.NewInt(x))
	yF := NewFieldElement(big.NewInt(223), big.NewInt(y))
	aF := NewFieldElement(big.NewInt(223), big.NewInt(a))
	bF := NewFieldElement(big.NewInt(223), big.NewInt(b))
	var p *Point
	assert.NotPanics(t, func() {
		p = NewEllipticCurvePoint(xF, yF, aF, bF)
	})
	return p
}

func TestPublicKeyAddress(t *testing.T) {
	/*
		1. if the address for the wallet is on the mainnet, set first byte to 0x00
		if testnet, set first byte to 0x6f

		2. Encode the public into SEC(compressed or uncompressed), do sha256 hash
		do ripemd160 hash, operation hash160

		3. combine the byte from step 1 with bytes array from step 2

		4. do a hash256 on the result of step 3, the get the first 4 bytes from the hash,
		append the 4 bytes at the end of bytes array from step 3, operation base58checksum

		5. combine bytes array from step 3 and step 4, theencode it by base58
	*/
	privateKey := NewPrivateKey(big.NewInt(5002))
	publicKey := privateKey.GetPublicKey()
	fmt.Printf("wallet address for 5002*G is %s\n", publicKey.Address(false, true))

	// 2020^5*G
	var expOp big.Int
	privateKey = NewPrivateKey(expOp.Exp(big.NewInt(2020), big.NewInt(5), nil))
	publicKey = privateKey.GetPublicKey()
	fmt.Printf("wallet address for 2020^5*G is %s\n", publicKey.Address(true, true))

	// 0x12345deadbeef
	p := new(big.Int)
	p.SetString("12345deadbeef", 16)
	privateKey = NewPrivateKey(p)
	publicKey = privateKey.GetPublicKey()
	fmt.Printf("wallet address for 0x12345deadbeef*G is %s\n", publicKey.Address(true, false))
}
