package main

import (
	"fmt"
	"time"

	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/util"
)

const ProtocolVersionByteSize = 2
const IndexByteSize = 2
const AmountByteSize = 8

type Transaction struct {
	ProtocolVersion uint16              `json:"protocolVersion"`
	Inputs          []TransactionInput  `json:"inputs"`
	Outputs         []TransactionOutput `json:"outputs"`
	Timestamp       time.Time           `json:"time"`
}

type TransactionInput struct {
	PreviousTransactionHash  [crypto.HashLength]byte      `json:"previousTransactionHash"`
	PreviousTransactionIndex uint16                       `json:"previousTransactionIndex"`
	SenderSignature          [crypto.SignatureLength]byte `json:"senderSignature"`
}

type TransactionOutput struct {
	ReceiverAddress [AddressLength]byte `json:"receiverAddress"`
	Amount          float64             `json:"amount"`
}

func NewCoinbaseTransaction(address [AddressLength]byte, amount float64) Transaction {
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

// Takes a transaction and returns a byte array representing the transaction
// TODO: Make the byte array conversion more efficient by preallocation
func (t Transaction) TransactionToByteArray() []byte {
	transactionBytes := make([]byte, 0)

	versionBytes := util.Uint16ToBytes(t.ProtocolVersion)
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

// Converts a TransactionInput into a byte array
func (t TransactionInput) TransactionInputToByteArray() []byte {
	inputBytes := make([]byte, 0)

	inputBytes = append(inputBytes, t.PreviousTransactionHash[:]...)

	indexBytes := util.Uint16ToBytes(t.PreviousTransactionIndex)
	inputBytes = append(inputBytes, indexBytes...)

	inputBytes = append(inputBytes, t.SenderSignature[:]...)

	return inputBytes
}

// Coverts a byte array into a TransactionInput
// TODO: Check safety of inputs
func ByteArrayToTransactionInput(bytes []byte) TransactionInput {
	hashBytes := bytes[:crypto.HashLength]
	indexBytes := bytes[crypto.HashLength : crypto.HashLength+IndexByteSize]
	signatureBytes := bytes[crypto.HashLength+IndexByteSize:]

	var hash [crypto.HashLength]byte
	var index uint16
	var signature [crypto.SignatureLength]byte

	copy(hash[:], hashBytes)
	index = util.BytesToUint16(indexBytes)
	copy(signature[:], signatureBytes)

	input := TransactionInput{
		PreviousTransactionHash:  hash,
		PreviousTransactionIndex: index,
		SenderSignature:          signature,
	}

	return input
}

func (t TransactionOutput) TransactionOutputToByteArray() []byte {
	outputBytes := make([]byte, 0)

	outputBytes = append(outputBytes, t.ReceiverAddress[:]...)

	amountBytes := util.Float64ToBytes(t.Amount)
	outputBytes = append(outputBytes, amountBytes...)

	return outputBytes
}

// Converts a byte array into a TransactionOutput
// TODO: Check safety of inputs
func ByteArrayToTransactionOutput(bytes []byte) TransactionOutput {
	addressBytes := bytes[:AddressLength]
	amountBytes := bytes[AddressLength:]

	var address [AddressLength]byte
	var amount float64

	copy(address[:], addressBytes)
	amount = util.BytesToFloat64(amountBytes)

	output := TransactionOutput{
		ReceiverAddress: address,
		Amount:          amount,
	}

	return output
}
