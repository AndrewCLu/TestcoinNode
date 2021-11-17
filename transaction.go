package main

import "fmt"

type Transaction struct {
	ProtocolVersion int                 `json:"protocolVersion"`
	Inputs          []transactionInput  `json:"inputs"`
	Outputs         []transactionOutput `json:"outputs"`
}

type transactionInput struct {
	PreviousTransactionHash  string `json:"previousTransactionHash"`
	PreviousTransactionIndex int    `json:"previousTransactionIndex"`
	SenderSignature          string `json:"senderSignature"`
}

type transactionOutput struct {
	ReceiverAddress string  `json:"receiverAddress"`
	Amount          float64 `json:"amount"`
}

func NewCoinbaseTransaction(address string, amount float64) Transaction {
	output := transactionOutput{ReceiverAddress: address, Amount: amount}
	transaction := Transaction{
		ProtocolVersion: CurrentProtocolVersion,
		Inputs:          []transactionInput{},
		Outputs:         []transactionOutput{output},
	}

	fmt.Printf("Coinbase transaction sending %v to %v\n", amount, address)
	return transaction
}
