package transaction

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
)

type Transaction struct {
	version   *big.Int
	txInputs  []*TransactionInput
	txOutputs []*TransactionOutput
	lockTime  *big.Int
	testnet   bool
}

func ParseTransaction(binary []byte) *Transaction {
	transaction := Transaction{}
	reader := bytes.NewReader(binary)
	bufReader := bufio.NewReader(reader)

	verBuf := make([]byte, 4)
	bufReader.Read(verBuf)

	version := LittleEndianToBigInt(verBuf, LITTLE_ENDIAN_4_BYTES)
	fmt.Printf("transaction version: %x\n", version)

	inputs := getInputCount(bufReader)
	transactionInputs := []*TransactionInput{}
	for i := 0; i < int(inputs.Int64()); i++ {
		input := NewTransactionInput(bufReader)
		transactionInputs = append(transactionInputs, input)
	}
	transaction.txInputs = transactionInputs

	return nil
}

func getInputCount(bufReader *bufio.Reader) *big.Int {
	/*
		if the first byte of input is 0, then witness transaction,
		we need to skip the first two bytes(0x00, 0x01)
	*/
	firstByte, err := bufReader.Peek(1)
	if err != nil {
		panic(err)
	}
	if firstByte[0] == 0x00 {
		// skip the first two bytes
		skipBuf := make([]byte, 2)
		_, err = bufReader.Read(skipBuf)
		if err != nil {
			panic(err)
		}
	}

	count := ReadVarint(bufReader)
	fmt.Printf("input count is: %x\n", count)
	return count
}
