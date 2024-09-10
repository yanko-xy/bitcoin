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
	fetcher                  *TransactionFetcher
}

func NewTransactionInput(reader *bufio.Reader) *TransactionInput {
	// first 32 bytes are hash256 of previous transaction
	transactionInput := &TransactionInput{}
	transactionInput.fetcher = NewTransactionInputFetch()

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
	scriptBuf := transactionInput.scriptSig.Serialize()
	fmt.Printf("script bytes: %x\n", scriptBuf)

	// last 4 bytes for sequence
	seqBytes := make([]byte, 4)
	reader.Read(seqBytes)
	transactionInput.sqeuence = LittleEndianToBigInt(seqBytes, LITTLE_ENDIAN_4_BYTES)

	return transactionInput
}

func (t *TransactionInput) Value(testnet bool) *big.Int {
	previousTxID := fmt.Sprintf("%x", t.previousTransactionID)
	previousTx := t.fetcher.Fetch(previousTxID, testnet)
	tx := ParseTransaction(previousTx)

	return &tx.txOutputs[t.previousTransactionIndex.Int64()].amount
}

func reverseByteSlice(bytes []byte) []byte {
	reverseBytes := []byte{}
	for i := len(bytes) - 1; i >= 0; i-- {
		reverseBytes = append(reverseBytes, bytes[i])
	}
	return reverseBytes
}
