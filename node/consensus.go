package node

import (
	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/block"
	"github.com/AndrewCLu/TestcoinNode/chain"
	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/transaction"
	"github.com/AndrewCLu/TestcoinNode/util"
)

// Returns if a transaction is valid or not based on the state of the ledger
// TODO: Make verifications of utxos existing in ledger more efficient
// TODO: Safety checks
func ValidateTransaction(tx transaction.Transaction) bool {
	var inputTotal uint64 = 0
	for _, input := range tx.Inputs {
		ptr := input.OutputPointer
		verification := input.Verification
		signature := verification.Signature
		senderPublicKey := verification.EncodedPublicKey
		senderAddress := account.GetAddressFromPublicKey(senderPublicKey)

		// Verify that the input is actually signed by the utxo possessor
		if !VerifyInput(senderPublicKey, ptr, signature) {
			return false
		}

		// Gets a list of valid utxos for the sender
		// Because we verified the input using the sender key,
		// any utxo in this list belongs to the account that signed the current input
		outputPointers, _ := chain.GetUnspentTransactions(senderAddress)

		// Find if a current utxo matches the one implied by the transaction
		match := false
		for _, comparePtr := range outputPointers {
			if ptr.Equal(comparePtr) {
				match = true
			}
		}

		if !match {
			return false
		}

		amount, _ := chain.GetOutputAmount(ptr)
		inputTotal += amount
	}

	var outputTotal uint64 = 0
	for _, output := range tx.Outputs {
		outputTotal += output.Amount
	}

	if len(tx.Inputs) != 0 && inputTotal != outputTotal {
		return false
	}

	return true
}

// Returns if a block is valid or not given the state of the ledger
func ValidateBlock(block block.Block) bool {
	return true
}

// Signs a transaction input
func SignInput(privateKey []byte,
	outputPointer transaction.TransactionOutputPointer,
) (signature [transaction.TransactionSignatureLength]byte) {
	hash := outputPointer.TransactionHash
	index := outputPointer.OutputIndex

	inputBytes := append(hash[:], util.Uint16ToBytes(index)...)

	signature, _ = crypto.SignByteArray(inputBytes, privateKey)

	return signature
}

// Verifies the signature of a transaction input
func VerifyInput(publicKey []byte,
	outputPointer transaction.TransactionOutputPointer,
	signature [transaction.TransactionSignatureLength]byte,
) (verified bool) {
	hash := outputPointer.TransactionHash
	index := outputPointer.OutputIndex

	inputBytes := append(hash[:], util.Uint16ToBytes(index)...)

	verified, _ = crypto.VerifyByteArray(inputBytes, publicKey, signature)

	return verified
}
