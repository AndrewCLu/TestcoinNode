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
const TransactionVerificationOffset = TransactionHashLength + TransactionIndexLength
const TransactionVerificationLengthLength = 2
const TransactionSignatureLength = crypto.SignatureLength

const TransactionAmountLength = 8
const TransactionOutputLength = protocol.AddressLength + TransactionAmountLength

type Transaction struct {
	ProtocolVersion uint16              `json:"protocolVersion"`
	Inputs          []TransactionInput  `json:"inputs"`
	Outputs         []TransactionOutput `json:"outputs"`
	Timestamp       time.Time           `json:"time"`
}

type TransactionInput struct {
	PreviousTransactionHash  [TransactionHashLength]byte  `json:"previousTransactionHash"`
	PreviousTransactionIndex uint16                       `json:"previousTransactionIndex"`
	VerificationLength       uint16                       `json:"verificationLength"`
	Verification             TransactionInputVerification `json:"verification"`
}

type TransactionInputVerification struct {
	Signature        [TransactionSignatureLength]byte `json:"signature"`
	EncodedPublicKey []byte                           `json:"encodedPublicKey"`
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
	senderPublicKey []byte,
	senderPrivateKey []byte,
	utxos []UnspentTransactionOutput,
	outputs []TransactionOutput,
) (t Transaction, success bool) {

	inputs := []TransactionInput{}
	var inputTotal uint64 = 0
	for _, utxo := range utxos {
		signature := SignInput(senderPrivateKey, utxo.TransactionHash, utxo.TransactionIndex)

		inputTotal += utxo.Amount

		verification := TransactionInputVerification{
			Signature:        signature,
			EncodedPublicKey: senderPublicKey,
		}

		input := TransactionInput{
			PreviousTransactionHash:  utxo.TransactionHash,
			PreviousTransactionIndex: utxo.TransactionIndex,
			VerificationLength:       uint16(len(verification.TransactionInputVerificationToByteArray())),
			Verification:             verification,
		}

		inputs = append(inputs, input)
	}

	var outputTotal uint64 = 0
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
		verificationOffset := currentByte + TransactionVerificationOffset
		verificationLength := util.BytesToUint16(bytes[verificationOffset : verificationOffset+TransactionVerificationLengthLength])

		inputLength := TransactionVerificationOffset + TransactionVerificationLengthLength + int(verificationLength)
		input := ByteArrayToTransactionInput(bytes[currentByte : currentByte+inputLength])
		inputs = append(inputs, input)
		currentByte += inputLength
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

	verificationLengthBytes := util.Uint16ToBytes(t.VerificationLength)
	inputBytes = append(inputBytes, verificationLengthBytes...)

	verificationBytes := t.Verification.TransactionInputVerificationToByteArray()
	inputBytes = append(inputBytes, verificationBytes...)

	return inputBytes
}

// Coverts a byte array into a TransactionInput
// TODO: Check safety of inputs
func ByteArrayToTransactionInput(bytes []byte) TransactionInput {
	currentByte := 0

	hashBytes := bytes[currentByte : currentByte+TransactionHashLength]
	currentByte += TransactionHashLength

	indexBytes := bytes[currentByte : currentByte+TransactionIndexLength]
	currentByte += TransactionIndexLength

	verificationLengthBytes := bytes[currentByte : currentByte+TransactionVerificationLengthLength]
	currentByte += TransactionVerificationLengthLength

	verificationBytes := bytes[currentByte:]

	var hash [TransactionHashLength]byte
	var index uint16
	var verificationLength uint16
	var verification TransactionInputVerification

	copy(hash[:], hashBytes)
	index = util.BytesToUint16(indexBytes)
	verificationLength = util.BytesToUint16(verificationLengthBytes)
	verification = ByteArrayToTransactionInputVerification(verificationBytes)

	input := TransactionInput{
		PreviousTransactionHash:  hash,
		PreviousTransactionIndex: index,
		VerificationLength:       verificationLength,
		Verification:             verification,
	}

	return input
}

func (t TransactionInputVerification) TransactionInputVerificationToByteArray() []byte {
	verificationBytes := make([]byte, 0)

	verificationBytes = append(verificationBytes, t.Signature[:]...)

	verificationBytes = append(verificationBytes, t.EncodedPublicKey...)

	return verificationBytes
}

func ByteArrayToTransactionInputVerification(bytes []byte) TransactionInputVerification {
	signatureBytes := bytes[:TransactionSignatureLength]
	publicKey := bytes[TransactionSignatureLength:]

	var signature [TransactionSignatureLength]byte
	copy(signature[:], signatureBytes)

	output := TransactionInputVerification{
		Signature:        signature,
		EncodedPublicKey: publicKey,
	}

	return output
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

// Signs a transaction input
func SignInput(privateKey []byte,
	hash [TransactionHashLength]byte,
	index uint16,
) (signature [TransactionSignatureLength]byte) {
	inputBytes := append(hash[:], util.Uint16ToBytes(index)...)

	signature, _ = crypto.SignByteArray(inputBytes, privateKey)

	return signature
}

// Verifies the signature of a transaction input
func VerifyInput(publicKey []byte,
	hash [TransactionHashLength]byte,
	index uint16, signature [TransactionSignatureLength]byte,
) (verified bool) {
	inputBytes := append(hash[:], util.Uint16ToBytes(index)...)

	verified, _ = crypto.VerifyByteArray(inputBytes, publicKey, signature)

	return verified
}

// Hashes a transaction
func (t Transaction) Hash() [TransactionHashLength]byte {
	bytes := t.TransactionToByteArray()
	return crypto.HashBytes(bytes)
}

// Checks if two transactions are equal
func (ta Transaction) Equal(tb Transaction) bool {
	return reflect.DeepEqual(ta.Hash(), tb.Hash())
}

// Checks if two unspent transaction outputs are equal
func (utxoa UnspentTransactionOutput) Equal(utxob UnspentTransactionOutput) bool {
	return reflect.DeepEqual(utxoa.TransactionHash, utxob.TransactionHash) &&
		utxoa.TransactionIndex == utxob.TransactionIndex &&
		reflect.DeepEqual(utxoa.ReceiverAddress, utxob.ReceiverAddress) &&
		utxoa.Amount == utxob.Amount
}
