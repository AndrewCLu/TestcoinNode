package pow

import (
	"bytes"
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/block"
	"github.com/AndrewCLu/TestcoinNode/chain"
	"github.com/AndrewCLu/TestcoinNode/common"
	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/protocol"
	"github.com/AndrewCLu/TestcoinNode/transaction"
	"github.com/AndrewCLu/TestcoinNode/util"
)

// Pow is a consensus mechanism based on proof-of-work
type Pow struct {
}

func New() (p *Pow, ok bool) {
	pow := Pow{}

	return &pow, true
}

// Returns if a transaction is valid or not based on the state of the ledger
func (pow *Pow) ValidatePendingTransaction(chain *chain.Chain, tx *transaction.Transaction) bool {
	var inputTotal uint64 = 0
	usedUtxoHashes := []common.Hash{}
	for _, input := range tx.Inputs {
		ptr := input.OutputPointer
		verification := input.Verification
		signature := verification.Signature
		senderPublicKey := verification.EncodedPublicKey
		senderAddress := account.GetAddressFromPublicKey(senderPublicKey)

		// Verify this utxo has not been used already
		utxoHash := ptr.Hash()
		for _, usedHash := range usedUtxoHashes {
			if utxoHash.Equal(usedHash) {
				return false
			}
		}
		usedUtxoHashes = append(usedUtxoHashes, utxoHash)

		// Verify that the input is actually signed by the utxo possessor
		if !pow.VerifyInput(senderPublicKey, ptr, signature) {
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

	return inputTotal > outputTotal
}

// Returns boolean indicating if a coinbase transaction is valid or not based onn the state of the ledger
func (pow *Pow) ValidateCoinbaseTransaction(chain *chain.Chain, coinbase *transaction.Transaction) bool {
	if len(coinbase.Inputs) > 0 {
		fmt.Println("Coinbase transaction cannot have any inputs")
		return false
	}

	if len(coinbase.Outputs) != 1 {
		fmt.Println("Coinbase transaction can only have one output")
	}

	_, lastBlockNum, lastBlockOk := chain.GetLastBlockInfo()
	if !lastBlockOk {
		fmt.Println("Could not get last block from current chain")
		return false
	}
	blockNum := lastBlockNum + 1
	blockReward := protocol.ComputeBlockReward(blockNum)
	if coinbase.Outputs[0].Amount != blockReward {
		fmt.Println("Incorrect block reward for given block number")
		return false
	}

	return true
}

// Returns if a block is valid or not given the state of the ledger
// TODO: Make sure input validation is done dynamically - each transaction should update the utxo state
func (pow *Pow) ValidateBlock(chn *chain.Chain, block *block.Block) bool {
	header := block.Header
	transactions := block.Body
	prevHash, prevBlockNum, success := chn.GetLastBlockInfo()
	// Previous block is retrievable
	if !success {
		return false
	}

	// Number of transactions does not exceed limit
	if len(transactions) > protocol.MaxTransactionsInBlock {
		return false
	}

	// PreviousBlockHash corresponds to last block
	if bytes.Compare(prevHash[:], header.PreviousBlockHash[:]) != 0 {
		return false
	}

	// Validate coinbase transaction
	coinbase := block.Coinbase
	if !pow.ValidateCoinbaseTransaction(chn, coinbase) {
		return false
	}

	// Creates a new chain to test validity of this block
	tempChain := chn.UnsafeCopy()
	transactionHashes := make([][]byte, len(transactions)+1)
	for i, tx := range transactions {
		if !pow.ValidatePendingTransaction(tempChain, tx) {
			return false
		}
		tempChain.AddTransaction(tx)
		hash := tx.Hash()
		transactionHashes[i] = hash[:]
	}
	transactionHashes[len(transactions)] = coinbase.Hash().Bytes()
	allTransactionsHash := crypto.HashBytes(util.ConcatByteSlices(transactionHashes))
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
func (pow *Pow) SignInput(privateKey []byte,
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
func (pow *Pow) VerifyInput(publicKey []byte,
	outputPointer *transaction.TransactionOutputPointer,
	signature *crypto.ECDSASignature,
) (verified bool) {
	hash := outputPointer.TransactionHash
	index := outputPointer.OutputIndex

	inputBytes := append(hash[:], util.Uint16ToBytes(index)...)

	verified = crypto.VerifyByteArray(inputBytes, publicKey, signature)

	return verified
}
