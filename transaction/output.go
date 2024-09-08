package transaction

import (
	"bufio"
	"math/big"
)

type TransactionOutput struct {
	// satoshi
	amount       big.Int
	scriptPubKey *ScriptSig
	reader       *bufio.Reader
}

func NewTransactionOutput(reader *bufio.Reader) *TransactionOutput {
	return &TransactionOutput{
		reader: reader,
	}
}
