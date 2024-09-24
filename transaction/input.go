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
	sequence                 *big.Int
	fetcher                  *TransactionFetcher
}

func InitTransactionInput(previousTx []byte, previousIndex *big.Int) *TransactionInput {
	return &TransactionInput{
		previousTransactionID:    previousTx,
		previousTransactionIndex: previousIndex,
		scriptSig:                nil,
		sequence:                 big.NewInt(int64(0xffffffff)),
	}
}

func (t *TransactionInput) String() string {
	return fmt.Sprintf("previous transaction: %x\n previous tx index: %x\n",
		t.previousTransactionID,
		t.previousTransactionIndex,
	)
}

func (t *TransactionInput) SetScript(sig *ScriptSig) {
	t.scriptSig = sig
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

	// 4 bytes for previous transaction index
	idx := make([]byte, 4)
	reader.Read(idx)
	transactionInput.previousTransactionIndex = LittleEndianToBigInt(idx, LITTLE_ENDIAN_4_BYTES)

	transactionInput.scriptSig = NewScriptSig(reader)

	// last 4 bytes for sequence
	seqBytes := make([]byte, 4)
	reader.Read(seqBytes)
	transactionInput.sequence = LittleEndianToBigInt(seqBytes, LITTLE_ENDIAN_4_BYTES)

	return transactionInput
}

func (t *TransactionInput) getPreviousTx(testnet bool) *Transaction {
	previousTxID := fmt.Sprintf("%x", t.previousTransactionID)
	previousTx := t.fetcher.Fetch(previousTxID, testnet)
	tx := ParseTransaction(previousTx)
	return tx
}

func (t *TransactionInput) Value(testnet bool) *big.Int {
	tx := t.getPreviousTx(testnet)
	return tx.txOutputs[t.previousTransactionIndex.Int64()].amount
}

func (t *TransactionInput) Script(testnet bool) *ScriptSig {
	previousTxID := fmt.Sprintf("%x", t.previousTransactionID)
	previousTx := t.fetcher.Fetch(previousTxID, testnet)
	tx := ParseTransaction(previousTx)

	scriptPubKey := tx.txOutputs[t.previousTransactionIndex.Int64()].scriptPubKey
	return t.scriptSig.Add(scriptPubKey)
}

func (t *TransactionInput) scriptPubKey(testnet bool) *ScriptSig {
	tx := t.getPreviousTx(testnet)
	return tx.txOutputs[t.previousTransactionIndex.Int64()].scriptPubKey
}

func (t *TransactionInput) ReplaceWithScriptPubKey(testnet bool) {
	t.scriptSig = t.scriptPubKey(testnet)
	fmt.Printf("scriptpubkey: %v\n", t.scriptSig)
}

func (t *TransactionInput) Serialize() []byte {
	result := make([]byte, 0)
	result = append(result, reverseByteSlice(t.previousTransactionID)...)
	result = append(result, BigIntToLittleEndian(t.previousTransactionIndex, LITTLE_ENDIAN_4_BYTES)...)
	result = append(result, t.scriptSig.Serialize()...)
	result = append(result, BigIntToLittleEndian(t.sequence, LITTLE_ENDIAN_4_BYTES)...)
	return result
}

func reverseByteSlice(bytes []byte) []byte {
	reverseBytes := []byte{}
	for i := len(bytes) - 1; i >= 0; i-- {
		reverseBytes = append(reverseBytes, bytes[i])
	}
	return reverseBytes
}
