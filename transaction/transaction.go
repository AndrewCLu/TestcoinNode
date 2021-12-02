package transaction

import (
	"fmt"
	"time"

	"github.com/AndrewCLu/TestcoinNode/common"
	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/protocol"
	"github.com/AndrewCLu/TestcoinNode/util"
)

const (
	NumInputOutputLength                = 2                                          // The number of bytes used to designate the number of inputs or outputs
	TransactionIndexLength              = 2                                          // The number of bytes used to designate the index of an output
	TransactionOutputPointerLength      = common.HashLength + TransactionIndexLength // The number of bytes in an output pointer
	TransactionVerificationLengthLength = 2                                          // The number of bytes used to designate the length of a TransactionVerification
	TransactionSignatureLengthLength    = 2                                          // The number of bytes used to designate a transaction signature length

	TransactionAmountLength = 8                                              // The number of bytes used to designate the transaction amount
	TransactionOutputLength = common.AddressLength + TransactionAmountLength // The nunmber of bytes in an output
)

// A transaction is a collection of inputs and outputs that sends Testcoin between addresses
type Transaction struct {
	ProtocolVersion uint16              `json:"protocolVersion"`
	Inputs          []TransactionInput  `json:"inputs"`
	Outputs         []TransactionOutput `json:"outputs"`
	Timestamp       time.Time           `json:"timestamp"`
}

// A transaction input contains a pointer to a previous trasnaction output and a verification for proving ownership
type TransactionInput struct {
	OutputPointer      TransactionOutputPointer     `json:"outputPointer"`
	VerificationLength uint16                       `json:"verificationLength"`
	Verification       TransactionInputVerification `json:"verification"`
}

// A transaction output pointer points to a previous transaction output
type TransactionOutputPointer struct {
	TransactionHash common.Hash `json:"transactionHash"`
	OutputIndex     uint16      `json:"outputIndex"`
}

// A transaction input verification provides proof that a transaction input is owned by the sender
type TransactionInputVerification struct {
	SignatureLength  uint16                `json:"signatureLength"`
	Signature        crypto.ECDSASignature `json:"signature"`
	EncodedPublicKey []byte                `json:"encodedPublicKey"`
}

// A transaction output designates some amount of Testcoin to go to a receiver address
type TransactionOutput struct {
	ReceiverAddress common.Address `json:"receiverAddress"`
	Amount          uint64         `json:"amount"`
}

// Generates a new transaction and returns a pointer to it
// Also returns boolean indicating success
func NewTransaction(
	inputs []TransactionInput,
	outputs []TransactionOutput,
) (t *Transaction, success bool) {

	transaction := Transaction{
		ProtocolVersion: protocol.CurrentProtocolVersion,
		Inputs:          inputs,
		Outputs:         outputs,
		Timestamp:       time.Now().Round(0),
	}

	return &transaction, true
}

// Takes a transaction and returns a byte array representing the transaction
func (t Transaction) Bytes() []byte {
	versionBytes := util.Uint16ToBytes(t.ProtocolVersion)

	numInputBytes := util.Uint16ToBytes(uint16(len(t.Inputs)))

	inputBytes := make([]byte, 0)
	for _, input := range t.Inputs {
		inputBytes = append(inputBytes, input.Bytes()...)
	}

	numOutputBytes := util.Uint16ToBytes(uint16(len(t.Outputs)))

	outputBytes := make([]byte, 0)
	for _, output := range t.Outputs {
		outputBytes = append(outputBytes, output.Bytes()...)
	}

	timeBytes, err := t.Timestamp.MarshalBinary()
	if err != nil {
		fmt.Printf("Error occurred creating byte array for transaction timestamp: %v\n", err)
		return []byte{}
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
// TODO: Error handling for out of bounds
func BytesToTransaction(bytes []byte) Transaction {
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
		input := BytesToTransactionInput(bytes[currentByte : currentByte+inputLength])
		inputs = append(inputs, input)
		currentByte += inputLength
	}

	numOutputs := int(util.BytesToUint16(bytes[currentByte : currentByte+NumInputOutputLength]))
	currentByte += NumInputOutputLength

	outputs := []TransactionOutput{}
	for i := 0; i < numOutputs; i += 1 {
		output := BytesToTransactionOutput(bytes[currentByte : currentByte+TransactionOutputLength])
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
func (t TransactionInput) Bytes() []byte {
	outputPointerBytes := t.OutputPointer.Bytes()

	verificationLengthBytes := util.Uint16ToBytes(t.VerificationLength)

	verificationBytes := t.Verification.Bytes()

	allBytes := [][]byte{
		outputPointerBytes,
		verificationLengthBytes,
		verificationBytes,
	}

	return util.ConcatByteSlices(allBytes)
}

// Coverts a byte array into a TransactionInput
// TODO: Error handling for out of bounds
func BytesToTransactionInput(bytes []byte) TransactionInput {
	currentByte := 0

	outputPointerBytes := bytes[currentByte : currentByte+TransactionOutputPointerLength]
	currentByte += TransactionOutputPointerLength

	verificationLengthBytes := bytes[currentByte : currentByte+TransactionVerificationLengthLength]
	currentByte += TransactionVerificationLengthLength

	verificationBytes := bytes[currentByte:]

	var outputPointer TransactionOutputPointer
	var verificationLength uint16
	var verification TransactionInputVerification

	outputPointer = BytesToTransactionOutputPointer(outputPointerBytes)
	verificationLength = util.BytesToUint16(verificationLengthBytes)
	verification = BytesToTransactionInputVerification(verificationBytes)

	input := TransactionInput{
		OutputPointer:      outputPointer,
		VerificationLength: verificationLength,
		Verification:       verification,
	}

	return input
}

// Converts a TransactionOutputPointer to bytes
func (ptr TransactionOutputPointer) Bytes() []byte {
	hashBytes := ptr.TransactionHash.Bytes()

	indexBytes := util.Uint16ToBytes(ptr.OutputIndex)

	allBytes := [][]byte{
		hashBytes,
		indexBytes,
	}

	return util.ConcatByteSlices(allBytes)
}

// Converts bytes back to a TransactionOutputPointer
// TODO: Error handling for out of bounds
func BytesToTransactionOutputPointer(bytes []byte) TransactionOutputPointer {
	hashBytes := bytes[:common.HashLength]
	indexBytes := bytes[common.HashLength:]

	hash := common.BytesToHash(hashBytes)
	index := util.BytesToUint16(indexBytes)

	ptr := TransactionOutputPointer{
		TransactionHash: hash,
		OutputIndex:     index,
	}

	return ptr
}

// Converts a TransactionInputVerification to bytes
func (t TransactionInputVerification) Bytes() []byte {
	signatureBytes := t.Signature.Bytes()

	publicKeyBytes := t.EncodedPublicKey

	allBytes := [][]byte{
		signatureBytes,
		publicKeyBytes,
	}

	return util.ConcatByteSlices(allBytes)
}

// Converts bytes to a TransactionInputVerification
// TODO: Error handling for out of bounds
func BytesToTransactionInputVerification(bytes []byte) TransactionInputVerification {
	currentByte := 0
	signatureLengthBytes := bytes[currentByte : currentByte+TransactionSignatureLengthLength]
	signatureLength := int(util.BytesToUint16(signatureLengthBytes))
	currentByte += TransactionSignatureLengthLength

	signatureBytes := bytes[currentByte : currentByte+signatureLength]
	signature := crypto.BytesToECDSASignature(signatureBytes)
	currentByte += signatureLength

	publicKey := bytes[currentByte:]

	output := TransactionInputVerification{
		Signature:        signature,
		EncodedPublicKey: publicKey,
	}

	return output
}

// Converts a TransactionOutput to bytes
func (t TransactionOutput) Bytes() []byte {
	addressBytes := t.ReceiverAddress.Bytes()

	amountBytes := util.Uint64ToBytes(t.Amount)

	allBytes := [][]byte{
		addressBytes,
		amountBytes,
	}

	return util.ConcatByteSlices(allBytes)
}

// Converts a byte array into a TransactionOutput
// TODO: Error handling for out of bounds
func BytesToTransactionOutput(bytes []byte) TransactionOutput {
	addressBytes := bytes[:common.AddressLength]
	amountBytes := bytes[common.AddressLength:]

	address := common.BytesToAddress(addressBytes)
	amount := util.BytesToUint64(amountBytes)

	output := TransactionOutput{
		ReceiverAddress: address,
		Amount:          amount,
	}

	return output
}

// Hashes a transaction
func (t Transaction) Hash() common.Hash {
	return crypto.HashBytes(t.Bytes())
}

// Checks if two transactions are equal
func (ta Transaction) Equal(tb Transaction) bool {
	hasha := ta.Hash()
	hashb := tb.Hash()

	return hasha.Equal(hashb)
}

// Checks if two unspent transaction outputs are equal
func (ptra TransactionOutputPointer) Equal(ptrb TransactionOutputPointer) bool {
	return ptra.TransactionHash.Equal(ptrb.TransactionHash) && ptra.OutputIndex == ptrb.OutputIndex
}
