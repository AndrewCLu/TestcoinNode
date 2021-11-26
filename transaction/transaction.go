package transaction

import (
	"fmt"
	"reflect"
	"time"

	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/protocol"
	"github.com/AndrewCLu/TestcoinNode/util"
)

const NumInputOutputLength = 2

const TransactionIndexLength = 2
const TransactionOutputPointerLength = crypto.HashLength + TransactionIndexLength
const TransactionVerificationLengthLength = 2
const TransactionSignatureLength = crypto.SignatureLength

const TransactionAmountLength = 8
const TransactionOutputLength = protocol.AddressLength + TransactionAmountLength

type Transaction struct {
	ProtocolVersion uint16              `json:"protocolVersion"`
	Inputs          []TransactionInput  `json:"inputs"`
	Outputs         []TransactionOutput `json:"outputs"`
	Timestamp       time.Time           `json:"timestamp"`
}

type TransactionInput struct {
	OutputPointer      TransactionOutputPointer     `json:"outputPointer"`
	VerificationLength uint16                       `json:"verificationLength"`
	Verification       TransactionInputVerification `json:"verification"`
}

type TransactionOutputPointer struct {
	TransactionHash [crypto.HashLength]byte `json:"transactionHash"`
	OutputIndex     uint16                  `json:"outputIndex"`
}

type TransactionInputVerification struct {
	Signature        [TransactionSignatureLength]byte `json:"signature"`
	EncodedPublicKey []byte                           `json:"encodedPublicKey"`
}

type TransactionOutput struct {
	ReceiverAddress [protocol.AddressLength]byte `json:"receiverAddress"`
	Amount          uint64                       `json:"amount"`
}

// Generates a new transaction and returns it
// Also returns boolean indicating success
func NewTransaction(
	inputs []TransactionInput,
	outputs []TransactionOutput,
) (t Transaction, success bool) {

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
	versionBytes := util.Uint16ToBytes(t.ProtocolVersion)

	numInputBytes := util.Uint16ToBytes(uint16(len(t.Inputs)))

	inputBytes := make([]byte, 0)
	for _, input := range t.Inputs {
		inputBytes = append(inputBytes, input.TransactionInputToByteArray()...)
	}

	numOutputBytes := util.Uint16ToBytes(uint16(len(t.Outputs)))

	outputBytes := make([]byte, 0)
	for _, output := range t.Outputs {
		outputBytes = append(outputBytes, output.TransactionOutputToByteArray()...)
	}

	timeBytes, err := t.Timestamp.MarshalBinary()
	if err != nil {
		fmt.Printf("Error occurred creating byte array for transaction timestamp: %v\n", err)
	}

	allBytes := [][]byte{
		versionBytes,
		numInputBytes,
		inputBytes,
		numOutputBytes,
		outputBytes,
		timeBytes,
	}

	return util.ConcatByteSlices(allBytes)
}

// Convertes a byte array back into a Transaction
// TODO: Check safety of inputs
func ByteArrayToTransaction(bytes []byte) Transaction {
	currentByte := 0

	protocolVersion := util.BytesToUint16(bytes[currentByte : currentByte+protocol.ProtocolVersionLength])
	currentByte += protocol.ProtocolVersionLength

	numInputs := int(util.BytesToUint16(bytes[currentByte : currentByte+NumInputOutputLength]))
	currentByte += NumInputOutputLength

	inputs := []TransactionInput{}
	for i := 0; i < numInputs; i += 1 {
		verificationOffset := currentByte + TransactionOutputPointerLength
		verificationLength := util.BytesToUint16(bytes[verificationOffset : verificationOffset+TransactionVerificationLengthLength])

		inputLength := TransactionOutputPointerLength + TransactionVerificationLengthLength + int(verificationLength)
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
	outputPointerBytes := t.OutputPointer.TransactionOutputPointerToByteArray()

	verificationLengthBytes := util.Uint16ToBytes(t.VerificationLength)

	verificationBytes := t.Verification.TransactionInputVerificationToByteArray()

	allBytes := [][]byte{
		outputPointerBytes,
		verificationLengthBytes,
		verificationBytes,
	}

	return util.ConcatByteSlices(allBytes)
}

// Coverts a byte array into a TransactionInput
// TODO: Check safety of inputs
func ByteArrayToTransactionInput(bytes []byte) TransactionInput {
	currentByte := 0

	outputPointerBytes := bytes[currentByte : currentByte+TransactionOutputPointerLength]
	currentByte += TransactionOutputPointerLength

	verificationLengthBytes := bytes[currentByte : currentByte+TransactionVerificationLengthLength]
	currentByte += TransactionVerificationLengthLength

	verificationBytes := bytes[currentByte:]

	var outputPointer TransactionOutputPointer
	var verificationLength uint16
	var verification TransactionInputVerification

	outputPointer = ByteArrayToTransactionOutputPointer(outputPointerBytes)
	verificationLength = util.BytesToUint16(verificationLengthBytes)
	verification = ByteArrayToTransactionInputVerification(verificationBytes)

	input := TransactionInput{
		OutputPointer:      outputPointer,
		VerificationLength: verificationLength,
		Verification:       verification,
	}

	return input
}

func (ptr TransactionOutputPointer) TransactionOutputPointerToByteArray() []byte {
	hashBytes := ptr.TransactionHash[:]

	indexBytes := util.Uint16ToBytes(ptr.OutputIndex)

	allBytes := [][]byte{
		hashBytes,
		indexBytes,
	}

	return util.ConcatByteSlices(allBytes)
}

func ByteArrayToTransactionOutputPointer(bytes []byte) TransactionOutputPointer {
	hashBytes := bytes[:crypto.HashLength]
	indexBytes := bytes[crypto.HashLength:]

	var hash [crypto.HashLength]byte
	copy(hash[:], hashBytes)
	index := util.BytesToUint16(indexBytes)

	ptr := TransactionOutputPointer{
		TransactionHash: hash,
		OutputIndex:     index,
	}

	return ptr
}

func (t TransactionInputVerification) TransactionInputVerificationToByteArray() []byte {
	signatureBytes := t.Signature[:]

	publicKeyBytes := t.EncodedPublicKey

	allBytes := [][]byte{
		signatureBytes,
		publicKeyBytes,
	}

	return util.ConcatByteSlices(allBytes)
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
	addressBytes := t.ReceiverAddress[:]

	amountBytes := util.Uint64ToBytes(t.Amount)

	allBytes := [][]byte{
		addressBytes,
		amountBytes,
	}

	return util.ConcatByteSlices(allBytes)
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

// Hashes a transaction
func (t Transaction) Hash() [crypto.HashLength]byte {
	return crypto.HashBytes(t.TransactionToByteArray())
}

// Checks if two transactions are equal
func (ta Transaction) Equal(tb Transaction) bool {
	return reflect.DeepEqual(ta.Hash(), tb.Hash())
}

// Checks if two unspent transaction outputs are equal
func (ptra TransactionOutputPointer) Equal(ptrb TransactionOutputPointer) bool {
	return reflect.DeepEqual(ptra.TransactionHash, ptrb.TransactionHash) && ptra.OutputIndex == ptrb.OutputIndex
}
