package transaction

import (
	"bufio"
	"fmt"
	"math/big"
)

type TransactionInput struct {
	previousTransactionID    []byte
	previousTransactionIndex *big.Int
	scriptSig                *ScriptSig
	sqeuence                 *big.Int
}

func NewTransactionInput(reader *bufio.Reader) *TransactionInput {
	// first 32 bytes are hash256 of previous transaction
	transactionInput := &TransactionInput{}
	previousTransaction := make([]byte, 32)
	reader.Read(previousTransaction)
	// convert it from little endian to big endian
	// reverse the byte array [0x01, 0x02, 0x03, 0x04] -> [0x04, 0x03, 0x02, 0x01]
	transactionInput.previousTransactionID = reverseByteSlice(previousTransaction)
	fmt.Printf("previous transaction id: %x\n", transactionInput.previousTransactionID)

	// 4 bytes for previous transaction index
	idx := make([]byte, 4)
	reader.Read(idx)
	transactionInput.previousTransactionIndex = LittleEndianToBigInt(idx, LITTLE_ENDIAN_4_BYTES)
	fmt.Printf("previous tansaction index: %x\n", transactionInput.previousTransactionIndex)

	transactionInput.scriptSig = NewScriptSig(reader)

	// last 4 bytes for sequence
	seqBytes := make([]byte, 4)
	reader.Read(seqBytes)
	transactionInput.sqeuence = LittleEndianToBigInt(seqBytes, LITTLE_ENDIAN_4_BYTES)

	return transactionInput
}

func reverseByteSlice(bytes []byte) []byte {
	reverseBytes := []byte{}
	for i := len(bytes) - 1; i >= 0; i-- {
		reverseBytes = append(reverseBytes, bytes[i])
	}
	return reverseBytes
}
