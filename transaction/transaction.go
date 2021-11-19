package transaction

import (
	"fmt"
	"reflect"
	"time"

	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/protocol"
	"github.com/AndrewCLu/TestcoinNode/util"
)

const ProtocolVersionLength = 2
const NumInputOutputLength = 2

const TransactionHashLength = crypto.HashLength
const TransactionIndexLength = 2
const TransactionSignatureLength = crypto.SignatureLength
const TransactionInputLength = TransactionHashLength + TransactionIndexLength + TransactionSignatureLength

const TransactionAmountLength = 8
const TransactionOutputLength = protocol.AddressLength + TransactionAmountLength

type Transaction struct {
	ProtocolVersion uint16              `json:"protocolVersion"`
	Inputs          []TransactionInput  `json:"inputs"`
	Outputs         []TransactionOutput `json:"outputs"`
	Timestamp       time.Time           `json:"time"`
}

type TransactionInput struct {
	PreviousTransactionHash  [TransactionHashLength]byte      `json:"previousTransactionHash"`
	PreviousTransactionIndex uint16                           `json:"previousTransactionIndex"`
	SenderSignature          [TransactionSignatureLength]byte `json:"senderSignature"`
}

type TransactionOutput struct {
	ReceiverAddress [protocol.AddressLength]byte `json:"receiverAddress"`
	Amount          uint64                       `json:"amount"`
}

type UnspentTransactionOutput struct {
	TransactionHash  [TransactionHashLength]byte  `json:"transactionHash"`
	TransactionIndex uint16                       `json:"transactionIndex"`
	ReceiverAddress  [protocol.AddressLength]byte `json:"receiverAddress"`
	Amount           uint64                       `json:"amount"`
}

// Generates a new coinbase transaction and returns it.
// Also returns boolean indicating success
func NewCoinbaseTransaction(
	address [protocol.AddressLength]byte,
	amount uint64) (t Transaction, success bool) {
	output := TransactionOutput{ReceiverAddress: address, Amount: amount}
	transaction := Transaction{
		ProtocolVersion: protocol.CurrentProtocolVersion,
		Inputs:          []TransactionInput{},
		Outputs:         []TransactionOutput{output},
		Timestamp:       time.Now().Round(0),
	}

	return transaction, true
}

// Generates a new peer transaction and returns it.
// Also returns boolean indicating success
func NewPeerTransaction(
	senderPrivateKey []byte,
	utxos []UnspentTransactionOutput,
	outputs []TransactionOutput,
) (t Transaction, success bool) {

	inputs := []TransactionInput{}
	inputTotal := uint64(0)
	for _, utxo := range utxos {
		transactionInputBytes := append(utxo.TransactionHash[:], util.Uint16ToBytes(utxo.TransactionIndex)...)

		inputTotal += utxo.Amount

		signature, _ := crypto.SignByteArray(transactionInputBytes, senderPrivateKey)
		input := TransactionInput{
			PreviousTransactionHash:  utxo.TransactionHash,
			PreviousTransactionIndex: utxo.TransactionIndex,
			SenderSignature:          signature,
		}

		inputs = append(inputs, input)
	}

	outputTotal := uint64(0)
	for _, output := range outputs {
		outputTotal += output.Amount
	}

	// If inputs and outputs don't match, transaction failed
	if inputTotal != outputTotal {
		return Transaction{}, false
	}

	transaction := Transaction{
		ProtocolVersion: protocol.CurrentProtocolVersion,
		Inputs:          inputs,
		Outputs:         outputs,
		Timestamp:       time.Now().Round(0),
	}

	return transaction, true
}

// Hashes a transaction
func (t Transaction) GetTransactionHash() [TransactionHashLength]byte {
	bytes := t.TransactionToByteArray()
	return crypto.HashBytes(bytes)
}

// Takes a transaction and returns a byte array representing the transaction
// TODO: Make the byte array conversion more efficient by preallocation
func (t Transaction) TransactionToByteArray() []byte {
	transactionBytes := make([]byte, 0)

	versionBytes := util.Uint16ToBytes(t.ProtocolVersion)
	transactionBytes = append(transactionBytes, versionBytes...)

	numInputBytes := util.Uint16ToBytes(uint16(len(t.Inputs)))
	transactionBytes = append(transactionBytes, numInputBytes...)

	inputBytes := make([]byte, 0)
	for _, input := range t.Inputs {
		inputBytes = append(inputBytes, input.TransactionInputToByteArray()...)
	}
	transactionBytes = append(transactionBytes, inputBytes...)

	numOutputBytes := util.Uint16ToBytes(uint16(len(t.Outputs)))
	transactionBytes = append(transactionBytes, numOutputBytes...)

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

// Convertes a byte array back into a Transaction
// TODO: Check safety of inputs
func ByteArrayToTransaction(bytes []byte) Transaction {
	currentByte := 0

	protocolVersion := util.BytesToUint16(bytes[currentByte : currentByte+ProtocolVersionLength])
	currentByte += ProtocolVersionLength

	numInputs := int(util.BytesToUint16(bytes[currentByte : currentByte+NumInputOutputLength]))
	currentByte += NumInputOutputLength

	inputs := []TransactionInput{}
	for i := 0; i < numInputs; i += 1 {
		input := ByteArrayToTransactionInput(bytes[currentByte : currentByte+TransactionInputLength])
		inputs = append(inputs, input)
		currentByte += TransactionInputLength
	}

	numOutputs := int(util.BytesToUint16(bytes[currentByte : currentByte+NumInputOutputLength]))
	currentByte += NumInputOutputLength

	outputs := []TransactionOutput{}
	for i := 0; i < numOutputs; i += 1 {
		output := ByteArrayToTransactionOutput(bytes[currentByte : currentByte+TransactionOutputLength])
		outputs = append(outputs, output)
		currentByte += TransactionOutputLength
	}

	timestamp := new(time.Time)
	timestamp.UnmarshalBinary(bytes[currentByte:])

	return Transaction{
		ProtocolVersion: protocolVersion,
		Inputs:          inputs,
		Outputs:         outputs,
		Timestamp:       *timestamp,
	}
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
	hashBytes := bytes[:TransactionHashLength]
	indexBytes := bytes[TransactionHashLength : TransactionHashLength+TransactionIndexLength]
	signatureBytes := bytes[TransactionHashLength+TransactionIndexLength:]

	var hash [TransactionHashLength]byte
	var index uint16
	var signature [TransactionSignatureLength]byte

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

	amountBytes := util.Uint64ToBytes(t.Amount)
	outputBytes = append(outputBytes, amountBytes...)

	return outputBytes
}

// Converts a byte array into a TransactionOutput
// TODO: Check safety of inputs
func ByteArrayToTransactionOutput(bytes []byte) TransactionOutput {
	addressBytes := bytes[:protocol.AddressLength]
	amountBytes := bytes[protocol.AddressLength:]

	var address [protocol.AddressLength]byte
	var amount uint64

	copy(address[:], addressBytes)
	amount = util.BytesToUint64(amountBytes)

	output := TransactionOutput{
		ReceiverAddress: address,
		Amount:          amount,
	}

	return output
}

// Checks if two transactions are equal
func AreTransactionsEqual(ta Transaction, tb Transaction) bool {
	return reflect.DeepEqual(ta.GetTransactionHash(), tb.GetTransactionHash())
}
