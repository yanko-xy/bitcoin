package transaction

import (
	"bufio"
	"bytes"
	ecc "elliptic_curve"
	"fmt"
	"math/big"
)

const (
	SIGHASH_ALL = 1
)

type Transaction struct {
	version   *big.Int
	txInputs  []*TransactionInput
	txOutputs []*TransactionOutput
	lockTime  *big.Int
	testnet   bool
}

func InitTransaction(version *big.Int, txInputs []*TransactionInput, txOutputs []*TransactionOutput, lockTime *big.Int, testnet bool) *Transaction {
	return &Transaction{
		version:   version,
		txInputs:  txInputs,
		txOutputs: txOutputs,
		lockTime:  lockTime,
		testnet:   testnet,
	}
}

func (t *Transaction) String() string {
	txIns := ""
	for i := 0; i < len(t.txInputs); i++ {
		txIns += t.txInputs[i].String()
		txIns += "\n"
	}
	txOuts := ""
	for i := 0; i < len(t.txOutputs); i++ {
		txOuts += t.txOutputs[i].String()
		txOuts += "\n"
	}
	return fmt.Sprintf("tx version: %v\n transaction inputs: %x\n transaction outputs: %x\n lockTime: %v\n",
		t.version,
		txIns,
		txOuts,
		t.lockTime,
	)
}

func (t *Transaction) SerializeWithSign(inputIdx int) []byte {
	/*
		constract signature message for the giving input indicate by input index,
		we need to change the given scriptsig with the scriptpubkey from the
		output of previous transaction, and the do hash256 on the binary transaction
		data
	*/
	signBinary := make([]byte, 0)
	signBinary = append(signBinary, BigIntToLittleEndian(t.version, LITTLE_ENDIAN_4_BYTES)...)

	inputCount := big.NewInt(int64(len(t.txInputs)))
	signBinary = append(signBinary, EncodeVarint(inputCount)...)

	/*
		serialize inputs, need to replace the scriptSig of the given input
		to scriptPubKey of previous transaction
	*/
	for i := 0; i < len(t.txInputs); i++ {
		if i == inputIdx {
			t.txInputs[i].ReplaceWithScriptPubKey(t.testnet)
			signBinary = append(signBinary, t.txInputs[i].Serialize()...)
		} else {
			signBinary = append(signBinary, t.txInputs[i].Serialize()...)
		}
	}

	outputCount := big.NewInt(int64(len(t.txOutputs)))
	signBinary = append(signBinary, EncodeVarint(outputCount)...)
	for i := 0; i < len(t.txOutputs); i++ {
		signBinary = append(signBinary, t.txOutputs[i].Serialize()...)
	}

	signBinary = append(signBinary, BigIntToLittleEndian(t.lockTime, LITTLE_ENDIAN_4_BYTES)...)
	signBinary = append(signBinary, BigIntToLittleEndian(big.NewInt(int64(SIGHASH_ALL)),
		LITTLE_ENDIAN_4_BYTES)...)

	return signBinary
}

func (t *Transaction) SignHash(inputIdx int) []byte {
	signBinary := t.SerializeWithSign(inputIdx)
	h256 := ecc.Hash256(string(signBinary))
	return h256
}

func (t *Transaction) VerifyInput(inputIdx int) bool {
	verifyScript := t.GetScript(inputIdx, t.testnet)
	z := t.SignHash(inputIdx)
	return verifyScript.Evaluate(z)
}

func (t *Transaction) Verify() bool {
	/*
		1. verify fee
		2. verify each transaction input
	*/
	if t.Fee().Cmp(big.NewInt(int64(0))) < 0 {
		return false
	}

	for i := 0; i < len(t.txInputs); i++ {
		if !t.VerifyInput(i) {
			return false
		}
	}

	return true
}

func ParseTransaction(binary []byte) *Transaction {
	transaction := &Transaction{}
	reader := bytes.NewReader(binary)
	bufReader := bufio.NewReader(reader)

	verBuf := make([]byte, 4)
	bufReader.Read(verBuf)

	version := LittleEndianToBigInt(verBuf, LITTLE_ENDIAN_4_BYTES)
	fmt.Printf("transaction version: %x\n", version)
	transaction.version = version

	inputs := getInputCount(bufReader)
	transactionInputs := []*TransactionInput{}
	for i := 0; i < int(inputs.Int64()); i++ {
		input := NewTransactionInput(bufReader)
		transactionInputs = append(transactionInputs, input)
	}
	transaction.txInputs = transactionInputs

	/*
		0100000001{813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c
		6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e
		24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7
		f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c
		7b8138bd94bdd531d2e213bf016b278afeffffff}
		{02}{a135ef0100000000}1976a914bc
		3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4
		bc762dd5423e332166702cb75f40df79fea1288ac{19430600}

		outputcount: 0x02

		output:
		1. amount in satoshi 1/ 1,000,000,000 of bitcoin (8 bytes)
		a135ef0100000000

		2. ScriptPubKey => ScriptSig

		locktime: 19430600
	*/

	// read output connts
	outputs := ReadVarint(bufReader)
	transactionOutputs := []*TransactionOutput{}
	for i := 0; i < int(outputs.Int64()); i++ {
		output := NewTransactionOutput(bufReader)
		transactionOutputs = append(transactionOutputs, output)
	}
	transaction.txOutputs = transactionOutputs

	// get last four byte for lock time
	lockTimeBytes := make([]byte, 4)
	bufReader.Read(lockTimeBytes)
	transaction.lockTime = LittleEndianToBigInt(lockTimeBytes, LITTLE_ENDIAN_4_BYTES)

	return transaction
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

func (t *Transaction) GetScript(idx int, testnet bool) *ScriptSig {
	if idx < 0 || idx > len(t.txInputs) {
		panic("invalid idx for transaction input")
	}

	txInputs := t.txInputs[idx]
	return txInputs.Script(testnet)
}

func (t *Transaction) Fee() *big.Int {
	// amount of input - amount of output > 0
	inputSum := big.NewInt(int64(0))
	outputSum := big.NewInt(int64(0))

	for i := 0; i < len(t.txInputs); i++ {
		addOP := new(big.Int)
		value := t.txInputs[i].Value(t.testnet)
		inputSum = addOP.Add(inputSum, value)
	}

	for i := 0; i < len(t.txOutputs); i++ {
		addOp := new(big.Int)
		outputSum = addOp.Add(outputSum, t.txOutputs[i].amount)
	}

	opSub := new(big.Int)
	return opSub.Sub(inputSum, outputSum)
}
