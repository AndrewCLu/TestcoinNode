package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

type Transaction struct {
	ProtocolVersion uint16              `json:"protocolVersion"`
	Inputs          []TransactionInput  `json:"inputs"`
	Outputs         []TransactionOutput `json:"outputs"`
	Timestamp       time.Time           `json:"time"`
}

type TransactionInput struct {
	PreviousTransactionHash  [32]byte `json:"previousTransactionHash"`
	PreviousTransactionIndex uint16   `json:"previousTransactionIndex"`
	SenderSignature          []byte   `json:"senderSignature"`
}

type TransactionOutput struct {
	ReceiverAddress [32]byte `json:"receiverAddress"`
	Amount          float64  `json:"amount"`
}

func NewCoinbaseTransaction(address [32]byte, amount float64) Transaction {
	output := TransactionOutput{ReceiverAddress: address, Amount: amount}
	transaction := Transaction{
		ProtocolVersion: CurrentProtocolVersion,
		Inputs:          []TransactionInput{},
		Outputs:         []TransactionOutput{output},
		Timestamp:       time.Now(),
	}

	fmt.Printf("Coinbase transaction sending %v to %v\n", amount, address)
	return transaction
}

func (t TransactionInput) TransactionInputToByteArray() []byte {
	inputBytes := make([]byte, 0)

	inputBytes = append(inputBytes, t.PreviousTransactionHash[:]...)

	indexBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(indexBytes, t.PreviousTransactionIndex)
	inputBytes = append(inputBytes, indexBytes...)

	inputBytes = append(inputBytes, t.SenderSignature...)

	return inputBytes
}

func (t TransactionOutput) TransactionOutputToByteArray() []byte {
	outputBytes := make([]byte, 0)

	outputBytes = append(outputBytes, t.ReceiverAddress[:]...)

	amountBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(amountBytes, math.Float64bits(t.Amount))
	outputBytes = append(outputBytes, amountBytes...)

	return outputBytes
}

// Takes a transaction and returns a byte array representing the transaction
// TODO: Make the byte array conversion more efficient by preallocation
func (t Transaction) TransactionToByteArray() []byte {
	transactionBytes := make([]byte, 0)

	versionBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(versionBytes, t.ProtocolVersion)
	transactionBytes = append(transactionBytes, versionBytes...)

	inputBytes := make([]byte, 0)
	for _, input := range t.Inputs {
		inputBytes = append(inputBytes, input.TransactionInputToByteArray()...)
	}
	transactionBytes = append(transactionBytes, inputBytes...)

	outputBytes := make([]byte, 0)
	for _, output := range t.Outputs {
		outputBytes = append(outputBytes, output.TransactionOutputToByteArray()...)
	}
	transactionBytes = append(transactionBytes, outputBytes...)

	timeBytes, err := t.Timestamp.MarshalBinary()
	if err != nil {
		fmt.Printf("Error occurred creating byte array for transaction timestamp: %v\n", err)
	}
	transactionBytes = append(transactionBytes, timeBytes...)

	return transactionBytes
}
