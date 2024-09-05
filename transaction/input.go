package transaction

import (
	"bufio"
)

type TransactionInput struct {
	reader *bufio.Reader
}

func NewTransactionInput(reader *bufio.Reader) *TransactionInput {
	return &TransactionInput{
		reader: reader,
	}
}
