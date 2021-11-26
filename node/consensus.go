package node

import (
	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/block"
	"github.com/AndrewCLu/TestcoinNode/chain"
	"github.com/AndrewCLu/TestcoinNode/transaction"
)

// Returns if a transaction is valid or not based on the state of the ledger
// TODO: Make verifications of utxos existing in ledger more efficient
// TODO: Safety checks
func ValidateTransaction(tx transaction.Transaction) bool {
	var inputTotal uint64 = 0
	for _, input := range tx.Inputs {
		hash := input.PreviousTransactionHash
		index := input.PreviousTransactionIndex
		verification := input.Verification
		signature := verification.Signature
		senderPublicKey := verification.EncodedPublicKey
		senderAddress := account.GetAddressFromPublicKey(senderPublicKey)

		// Make sure the input matches a real transaction
		previousTransaction, success := chain.GetTransaction(hash)
		if !success {
			return false
		}
		previousTransactionOutput := previousTransaction.Outputs[index]

		// Verify that the input is actually signed by the utxo possessor
		// TODO: Separate this function into the consensus package
		if !transaction.VerifyInput(senderPublicKey, hash, index, signature) {
			return false
		}

		// Make a utxo with the information from the previous transaction and see if it matches a valid utxo
		utxo := transaction.UnspentTransactionOutput{
			TransactionHash:  hash,
			TransactionIndex: index,
			ReceiverAddress:  previousTransactionOutput.ReceiverAddress,
			Amount:           previousTransactionOutput.Amount,
		}

		// Gets a list of valid utxos for the sender
		// Because we verified the input using the sender key,
		// any utxo in this list belongs to the account that signed the current input
		utxos, _ := chain.GetUnspentTransactions(senderAddress)

		// Find if a current utxo matches the one implied by the transaction
		match := false
		for _, compareUTXO := range utxos {
			if utxo.Equal(compareUTXO) {
				match = true
			}
		}

		if !match {
			return false
		}

		inputTotal += utxo.Amount
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
