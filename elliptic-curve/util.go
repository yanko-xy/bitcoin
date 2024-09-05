package elliptic_curve

import (
	"crypto/sha256"
	"math/big"

	"github.com/tsuna/endian"
	"golang.org/x/crypto/ripemd160"
)

/*
z, sha256(sha256(create a text)) -> 256bits -> 32bytes integer
*/
func Hash256(text string) []byte {
	hashOnce := sha256.Sum256([]byte(text))
	hashTwice := sha256.Sum256(hashOnce[:])
	return hashTwice[:]
}

func GetGenerator() *Point {
	Gx := new(big.Int)
	Gx.SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
	Gy := new(big.Int)
	Gy.SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)
	return S256Point(Gx, Gy)
}

func GetBitcoinValueN() *big.Int {
	n := new(big.Int)
	n.SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	return n
}

func ParseSEC(secBin []byte) *Point {
	// check the first byte to descide it is compressed or uncompressed
	if secBin[0] == 4 {
		// uncompressed
		x := new(big.Int)
		x.SetBytes(secBin[1:33])
		y := new(big.Int)
		y.SetBytes(secBin[33:65])
		return S256Point(x, y)
	}

	// check first byte for y is odd or even
	isEven := (secBin[0] == 2)
	x := new(big.Int)
	x.SetBytes(secBin[1:])
	y2 := S256Field(x).Power(big.NewInt(3)).Add(S256Field(big.NewInt(7)))
	y := y2.Sqrt()
	var modOp big.Int
	var yEven *FieldElement
	var yOdd *FieldElement
	if modOp.Mod(y.num, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		yEven = y
		yOdd = y.Negate() // p-y
	} else {
		yOdd = y
		yEven = y.Negate()
	}

	if isEven {
		return S256Point(x, yEven.num)
	} else {
		return S256Point(x, yOdd.num)
	}
}

/*
base58 it removes 0 o I l
*/
func EncodeBase58(s []byte) string {
	BASE58_ALPHABET := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	count := 0
	for idx := range s {
		if s[idx] == 0 {
			count += 1
		} else {
			break
		}
	}

	prefix := ""
	for i := 0; i < count; i++ {
		prefix += "1"
	}

	result := ""
	num := new(big.Int)
	num.SetBytes(s)
	for num.Cmp(big.NewInt(0)) > 0 {
		var divOp big.Int
		var modOp big.Int
		mod := modOp.Mod(num, big.NewInt(58))
		num = divOp.Div(num, big.NewInt(58))
		result = string(BASE58_ALPHABET[mod.Int64()]) + result
	}

	return prefix + result
}

func Base58Checksum(s []byte) string {
	hash256 := Hash256(string(s))
	return EncodeBase58(append(s, hash256[:4]...))
}

func Hash160(s []byte) []byte {
	sha256Bytes := sha256.Sum256(s)
	hasher := ripemd160.New()
	hasher.Write(sha256Bytes[:])
	hashBytes := hasher.Sum(nil)

	return hashBytes
}

type LITTLE_ENDIAN_LENGTH int

const (
	LITTLE_ENDIAN_2_BYTES = iota
	LITTLE_ENDIAN_4_BYTES
	LITTLE_ENDIAN_8_BYTES
)

func BigIntToLittleEndian(v *big.Int, length LITTLE_ENDIAN_LENGTH) []byte {
	switch length {
	case LITTLE_ENDIAN_2_BYTES:
		val := v.Int64()
		littleEndianVal := endian.HostToNetUint16(uint16(val))
		p := big.NewInt(int64(littleEndianVal))
		return p.Bytes()
	case LITTLE_ENDIAN_4_BYTES:
		val := v.Int64()
		littleEndianVal := endian.HostToNetUint32(uint32(val))
		p := big.NewInt(int64(littleEndianVal))
		return p.Bytes()
	case LITTLE_ENDIAN_8_BYTES:
		val := v.Int64()
		littleEndianVal := endian.HostToNetUint64(uint64(val))
		p := big.NewInt(int64(littleEndianVal))
		return p.Bytes()
	}

	return nil
}

func LittleEndianToBigInt(bytes []byte, length LITTLE_ENDIAN_LENGTH) *big.Int {
	switch length {
	case LITTLE_ENDIAN_2_BYTES:
		p := new(big.Int)
		p.SetBytes(bytes)
		val := endian.NetToHostUint16(uint16(p.Uint64()))
		return big.NewInt(int64(val))
	case LITTLE_ENDIAN_4_BYTES:
		p := new(big.Int)
		p.SetBytes(bytes)
		val := endian.NetToHostUint32(uint32(p.Uint64()))
		return big.NewInt(int64(val))
	case LITTLE_ENDIAN_8_BYTES:
		p := new(big.Int)
		p.SetBytes(bytes)
		val := endian.NetToHostUint64(uint64(p.Uint64()))
		return big.NewInt(int64(val))
	}

	return nil
}
