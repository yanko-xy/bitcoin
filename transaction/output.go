package transaction

import (
	"bufio"
)

type TransactionOutput struct {
	reader *bufio.Reader
}

func NewTransactionOutput(reader *bufio.Reader) *TransactionOutput {
	return &TransactionOutput{
		reader: reader,
	}
}
