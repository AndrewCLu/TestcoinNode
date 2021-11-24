package node

import (
	"github.com/AndrewCLu/TestcoinNode/account"
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

		previousTransaction, _ := chain.GetTransaction(hash)
		previousTransactionOutput := previousTransaction.Outputs[index]
		utxo := transaction.UnspentTransactionOutput{
			TransactionHash:  hash,
			TransactionIndex: index,
			ReceiverAddress:  previousTransactionOutput.ReceiverAddress,
			Amount:           previousTransactionOutput.Amount,
		}

		// Check if input is provided by the sender
		// TODO: Separate this function into the consensus package
		if !transaction.VerifyInput(senderPublicKey, hash, index, signature) {
			return false
		}

		// Find if a current utxo matches the one implied by the transaction
		match := false
		utxos, _ := chain.GetUnspentTransactions(senderAddress)
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
