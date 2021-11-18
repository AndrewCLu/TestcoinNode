package main

import (
	"encoding/binary"
	"fmt"
	"time"
)

type Transaction struct {
	ProtocolVersion uint16              `json:"protocolVersion"`
	Inputs          []TransactionInput  `json:"inputs"`
	Outputs         []TransactionOutput `json:"outputs"`
	Timestamp       time.Time           `json:"time"`
}

type TransactionInput struct {
	PreviousTransactionHash  string `json:"previousTransactionHash"`
	PreviousTransactionIndex int    `json:"previousTransactionIndex"`
	SenderSignature          string `json:"senderSignature"`
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
	return make([]byte, 2)
}

func (t TransactionOutput) TransactionOutputToByteArray() []byte {
	return make([]byte, 2)
}

// Takes a transaction and returns a byte array representing the transaction
// TODO: Make the byte array conversion more efficient by preallocation
func (t Transaction) TransactionToByteArray() []byte {
	transactionBytes := make([]byte, 0)

	versionBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(versionBytes, t.ProtocolVersion)
	fmt.Println(versionBytes)
	transactionBytes = append(transactionBytes, versionBytes...)

	inputBytes := make([]byte, 0)
	for _, input := range t.Inputs {
		inputBytes = append(inputBytes, input.TransactionInputToByteArray()...)
	}
	fmt.Println(inputBytes)
	transactionBytes = append(transactionBytes, inputBytes...)

	outputBytes := make([]byte, 0)
	for _, output := range t.Outputs {
		outputBytes = append(outputBytes, output.TransactionOutputToByteArray()...)
	}
	fmt.Println(outputBytes)
	transactionBytes = append(transactionBytes, outputBytes...)

	timeBytes, err := t.Timestamp.MarshalBinary()
	if err != nil {
		fmt.Printf("Error occurred creating byte array for transaction timestamp: %v\n", err)
	}
	fmt.Println(timeBytes)
	transactionBytes = append(transactionBytes, timeBytes...)

	return transactionBytes
}
