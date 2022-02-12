package consensus

import (
	"bytes"

	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/block"
	"github.com/AndrewCLu/TestcoinNode/chain"
	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/protocol"
	"github.com/AndrewCLu/TestcoinNode/transaction"
	"github.com/AndrewCLu/TestcoinNode/util"
)

// Returns if a transaction is valid or not based on the state of the ledger
// TODO: Check that pointers aren't reused as different inputs in the same tx
// TODO: Safety checks
func ValidateTransaction(chain *chain.Chain, tx *transaction.Transaction) bool {
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
// TODO: Make sure input validation is done dynamically - each transaction should update the utxo state
func ValidateBlock(chain *chain.Chain, block *block.Block) bool {
	header := block.Header
	transactions := block.Body
	prevHash, prevBlockNum, success := chain.GetLastBlockInfo()
	// Previous block is retrievable
	if !success {
		return false
	}

	// PreviousBlockHash corresponds to last block
	if bytes.Compare(prevHash[:], header.PreviousBlockHash[:]) != 0 {
		return false
	}

	transactionHashes := make([][]byte, len(transactions))
	for i, tx := range transactions {
		// Check that each transaction is valid
		// TODO: UPDATE STATE IN BETWEEN CHECKING IF TRANSACTIONS ARE VALID
		if !ValidateTransaction(chain, tx) {
			return false
		}

		hash := tx.Hash()
		transactionHashes[i] = hash[:]
	}
	allTransactionsHash := crypto.HashBytes(util.ConcatByteSlices(transactionHashes))
	// AllTransactionsHash is the hash of all transactionns
	if bytes.Compare(allTransactionsHash[:], header.AllTransactionsHash[:]) != 0 {
		return false
	}

	blockNum := prevBlockNum + 1
	targetHeader := protocol.ComputeTarget(blockNum)
	// Check that the selected target is correct
	if bytes.Compare(targetHeader[:], header.Target[:]) != 0 {
		return false
	}

	headerHash := header.Hash()
	target := targetHeader.FullHash()
	// Check that the computed hash is valid based on the target
	if bytes.Compare(headerHash[:], target[:]) >= 0 {
		return false
	}

	return true
}

// Signs a transaction input
// TODO: Should return a transaction input verification
func SignInput(privateKey []byte,
	outputPointer *transaction.TransactionOutputPointer,
) (signature *crypto.ECDSASignature) {
	hash := outputPointer.TransactionHash
	index := outputPointer.OutputIndex

	inputBytes := append(hash[:], util.Uint16ToBytes(index)...)

	signature, _ = crypto.SignByteArray(inputBytes, privateKey)

	return signature
}

// Verifies the signature of a transaction input
// TODO: Should accept a transaction input verification
func VerifyInput(publicKey []byte,
	outputPointer *transaction.TransactionOutputPointer,
	signature *crypto.ECDSASignature,
) (verified bool) {
	hash := outputPointer.TransactionHash
	index := outputPointer.OutputIndex

	inputBytes := append(hash[:], util.Uint16ToBytes(index)...)

	verified = crypto.VerifyByteArray(inputBytes, publicKey, signature)

	return verified
}
