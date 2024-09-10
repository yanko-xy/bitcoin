package transaction

type BitcoinOpCode struct {
}

func NewBitcoinOpCode() *BitcoinOpCode {
	return &BitcoinOpCode{}
}

func (b *BitcoinOpCode) EncodeNum(num int64) []byte {
	if num == 0 {
		// not push 0x00 but empty byte string
		return []byte("")
	}

	result := []byte{}
	absNum := num
	negative := false
	if num < 0 {
		absNum = -num
		negative = true
	}

	for absNum > 0 {
		/*
			append the last byte of absNum into result
			notices result will be little endian byte array of absNum
		*/
		result = append(result, byte(absNum&0xff))
		absNum >>= 8
	}

	/*
		check the most significant bit, notice the most significant byte is
		at the end of result
		0x8080 -> 32896 -32896
	*/
	if (result[len(result)-1] & 0x80) != 0 {
		if negative {
			// need to insert 0x80 at the head, most significant byte at the end
			// of result, we should insert 0x80 at the end
			result = append(result, 0x80)
		} else {
			result = append(result, 0x00)
		}
	} else if negative {
		// set the most significant byte to 1
		result[len(result)-1] |= 0x80
	}

	return result
}

func (b *BitcoinOpCode) DecodeNum(element []byte) int64 {
	bidEndian := reverseByteSlice(element)
	negative := false
	result := int64(0)

	// if the most significant bit is 1, it is negative value
	if (bidEndian[0] & 0x80) != 0 {
		negative = true
		// reset the most significant bit to 0
		// 0x7f is 0111 11111
		result = int64(bidEndian[0] & 0x7f)
	} else {
		negative = false
		result = int64(bidEndian[0])
	}

	for i := 1; i < len(bidEndian); i++ {
		result <<= 8
		result += int64(bidEndian[i])
	}

	if negative {
		result = -result
	}

	return result
}
