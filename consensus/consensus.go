package consensus

import (
	"github.com/AndrewCLu/TestcoinNode/block"
	"github.com/AndrewCLu/TestcoinNode/chain"
	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/transaction"
)

// A consensus provides agreed upon methods for coordinating a shared state of the blockchain
type Consensus interface {
	// Returns a boolean indicating if the given transaction is valid based on the current state of the blockchain
	ValidatePendingTransaction(chain *chain.Chain, tx *transaction.Transaction) bool

	// Returns a boolean indicating if the given block is valid based on the current state of the blockchain
	ValidateBlock(chain *chain.Chain, block *block.Block) bool

	// Given a private key and a transaction output pointer, returns a valid signature for the given output
	SignInput(privateKey []byte, outputPointer *transaction.TransactionOutputPointer) *crypto.ECDSASignature

	// Given a public key and a signature of a transaction output, returns a boolean indicating if the signature is valid
	VerifyInput(publicKey []byte, outputPointer *transaction.TransactionOutputPointer, signature *crypto.ECDSASignature) bool
}
