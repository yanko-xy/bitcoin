package transaction

import (
	"bufio"
	"math/big"
)

type ScriptSig struct {
	cmds          [][]byte
	bitcoinOpCode *BitcoinOpCode
}

const (
	// [0x01, 0x4d] -> [1, 75]
	SCRIPT_DATA_LENGTH_BEGIN = 1
	SCRIPT_DATA_LENGTH_END   = 75
	OP_PUSHDATA1             = 76
	OP_PUSHDATA2             = 77
)

func InitScriptSig(cmds [][]byte) *ScriptSig {
	bitcoinOpCode := NewBitcoinOpCode()
	bitcoinOpCode.cmds = cmds
	return &ScriptSig{
		bitcoinOpCode: bitcoinOpCode,
	}
}

/*
one kind for data operation -> move a chunk of data to stack
one kind for data processing -> get data of the top of stack and do something computation
and push the result on to the stack

mov eax, 0x1234

stack: [] 0x1234

it is not allowed for loop -> turing incomplete

byte value -> n, n [0x1, 0x4b] -> data operation, n is length of the chunk of data
we need to put on the top of stack
[0x4, 0x1, 0x2, 0x3, 0x4]
0x04 -> stack: [0x01020304]

4b -> 75, how can we move than 75 bytees on to the stack
n = 0x4c -> OP_PUSHDATA1, the following one byte is the length of the chunk of data
[0x4c, 0xfe, 0x1, ...]
how we read more than 75 bytes of data on to the stack
at most read 0xff 255, if we want read more than 255 bytes,
0x4d -> OP_PUSHDATA2, read following two bytes as the length of the chunk of data
[0x4d, 0x01, 0x02, ...]
0x0102 -> big endian -> 0x0201

OP_DUP https://en.bitcoin.it/wiki/Script
stack: [0x01020304, 0x01020304]

OP_ADD
stack: [0x01020304 + 0x01020304]

parsing script
read one byte, data operation, get length of data, move chunk of data on to stack data processing,
get elements from the stack, do the operation, push the result on to the stack
*/

func NewScriptSig(reader *bufio.Reader) *ScriptSig {
	cmds := [][]byte{}
	/*
		At the beginning is the total length for script field
	*/
	scriptLen := ReadVarint(reader).Int64()
	count := int64(0)
	current := make([]byte, 1)
	var current_byte byte
	for count < scriptLen {
		reader.Read(current)
		//operation
		count += 1
		current_byte = current[0]
		if current_byte >= SCRIPT_DATA_LENGTH_BEGIN &&
			current_byte <= SCRIPT_DATA_LENGTH_END {
			//push the following bytes of data onto stack
			data := make([]byte, current_byte)
			reader.Read(data)
			cmds = append(cmds, data)
			count += int64(current_byte)
		} else if current_byte == OP_PUSHDATA1 {
			/*
				read the following byte as the length of data
			*/
			length := make([]byte, 1)
			reader.Read(length)

			data := make([]byte, length[0])
			reader.Read(data)
			cmds = append(cmds, data)
			count += int64(length[0] + 1)
		} else if current_byte == OP_PUSHDATA2 {
			/*
				read the following 2 bytes as length of data
			*/
			lenBuf := make([]byte, 2)
			reader.Read(lenBuf)
			length := LittleEndianToBigInt(lenBuf, LITTLE_ENDIAN_2_BYTES)
			data := make([]byte, length.Int64())
			reader.Read(data)
			cmds = append(cmds, data)
			count += int64(2 + length.Int64())
		} else {
			//is data processing instruction
			cmds = append(cmds, []byte{current_byte})
		}
	}

	if count != scriptLen {
		panic("parsing script field fail")
	}

	return InitScriptSig(cmds)
}

func (s *ScriptSig) Evaluate(z []byte) bool {
	for s.bitcoinOpCode.HasCmd() {
		cmd := s.bitcoinOpCode.RemoveCmd()
		if len(cmd) == 1 {
			//this is an op code, run it
			opRes := s.bitcoinOpCode.ExecuteOperaion(int(cmd[0]), z)
			if !opRes {
				return false
			}
		} else {
			s.bitcoinOpCode.AppendDataElement(cmd)
		}
	}

	/*
		After runing all the operations in the scripts and the stack is empty
		then evaluation fail, otherwise we check the top element of the stack,
		if it value is 0, then fail, of the value is not 0, then success
	*/
	if len(s.bitcoinOpCode.stack) == 0 {
		return false
	}
	if len(s.bitcoinOpCode.stack[0]) == 0 {
		return false
	}

	return true
}

func (s *ScriptSig) rawSerialize() []byte {
	result := []byte{}
	for _, cmd := range s.bitcoinOpCode.cmds {
		if len(cmd) == 1 {
			//only one byte means its an instruction
			result = append(result, cmd...)
		} else {
			length := len(cmd)
			if length <= SCRIPT_DATA_LENGTH_END {
				//length in [0x01, 0x4b]
				result = append(result, byte(length))
			} else if length > SCRIPT_DATA_LENGTH_END && length < 0x100 {
				//this is OP_PUSHDATA1 command,
				//push the command and then the next byte is the length of the data
				result = append(result, OP_PUSHDATA1)
				result = append(result, byte(length))
			} else if length >= 0x100 && length <= 520 {
				/*
					this is OP_PUSHDATA2 command, we push the command
					and then two byte for the data length but in little endian format
				*/
				result = append(result, OP_PUSHDATA2)
				lenBuf := BigIntToLittleEndian(big.NewInt(int64(length)), LITTLE_ENDIAN_2_BYTES)
				result = append(result, lenBuf...)
			} else {
				panic("too long an cmd")
			}

			//append the chunk of data with given length
			result = append(result, cmd...)
		}
	}

	return result
}

func (s *ScriptSig) Serialize() []byte {
	rawResult := s.rawSerialize()
	total := len(rawResult)
	result := []byte{}
	// encode the total length of script at the head
	result = append(result, EncodeVarint(big.NewInt(int64(total)))...)
	result = append(result, rawResult...)
	return result
}

/*
1. move data around the stack, push some value onto the stack,
move value from one stack to another

OP_1, OP_2,..., OP_16, push one onto stack
OP_2 push value 2 on the stack

push value onto stack, transfer the value onto byte array and in little edian order
1234 -> 0x04de
-1234 use two's complement

1 -> 0x01 -> 0000 0001 -> -1 => 1000 0001 -> 0x81 -1 => [0x81]

-1234 => 12345 0x04 d2 => 1000 0100 1101 0010 => 0x84 d2 => [0x84, 0xd2]

0x8080 => 1000 0000 1000 0000 insert 0x80 at head
-32892 => 32892 => 0x8080 => 1000 0000 1000 0000 1000 0000 => 0x808080
=> [0x80, 0x80, 0x80] => -32892

*/

func (s *ScriptSig) Add(script *ScriptSig) *ScriptSig {
	cmds := make([][]byte, 0)
	cmds = append(cmds, s.bitcoinOpCode.cmds...)
	cmds = append(cmds, s.bitcoinOpCode.cmds...)
	return InitScriptSig(cmds)
}
