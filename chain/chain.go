package chain

import (
	"github.com/AndrewCLu/TestcoinNode/block"
	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/protocol"
	"github.com/AndrewCLu/TestcoinNode/transaction"
)

var blocks map[[crypto.HashLength]byte]block.Block                   // Stores all previous blocks
var transactions map[[crypto.HashLength]byte]transaction.Transaction // Stores all previous transactions

var lastBlockHash [crypto.HashLength]byte                                                  // Last block hash
var pendingTransactions []transaction.Transaction                                          // All transactions that have not been processed yet
var unspentOutputs map[[protocol.AddressLength]byte][]transaction.TransactionOutputPointer // All utxos

func InitializeChain() {
	genesisBlock, _ := block.NewBlock(crypto.HashBytes([]byte("first")), 0, []transaction.Transaction{})
	genesisBlockHash := genesisBlock.Hash()
	blocks = make(map[[crypto.HashLength]byte]block.Block)
	blocks[genesisBlockHash] = genesisBlock
	lastBlockHash = genesisBlockHash

	pendingTransactions = []transaction.Transaction{}
	transactions = make(map[[crypto.HashLength]byte]transaction.Transaction)
	unspentOutputs = make(map[[protocol.AddressLength]byte][]transaction.TransactionOutputPointer)
}

// Gets a recorded transaction
func GetTransaction(hash [crypto.HashLength]byte) (tx transaction.Transaction, success bool) {
	return transactions[hash], true
}

// Gets num pending transactions
// TODO: Have a ranking of pending transactions to retrieve by time added or miner fee
func GetPendingTransactions(num int) (txs []transaction.Transaction, success bool) {
	if num > len(pendingTransactions) {
		return pendingTransactions, true
	}

	return pendingTransactions[:num], true
}

// Add a pending transaction to the list
func AddPendingTransaction(tx transaction.Transaction) (success bool) {
	pendingTransactions = append(pendingTransactions, tx)

	return true
}

// Get information about the last block in the chain
func GetLastBlockInfo() (hash [crypto.HashLength]byte, blockNum int, success bool) {
	return lastBlockHash, len(blocks) - 1, true
}

// Adds a transaction to the chain
// TODO: Make this atomic - either it updates entire state if success or not at all
func AddTransaction(tx transaction.Transaction) (success bool) {
	// Find matching pending transaction
	ind := -1
	for i, ptx := range pendingTransactions {
		if tx.Equal(ptx) {
			ind = i
		}
	}

	// No pending transaction matches
	if ind == -1 {
		return false
	}

	// Remove pending transaction at index ind
	pendingTransactions[ind] = pendingTransactions[len(pendingTransactions)-1]
	pendingTransactions = pendingTransactions[:len(pendingTransactions)-1]

	// Add new transaction
	txHash := tx.Hash()
	transactions[txHash] = tx

	for _, input := range tx.Inputs {
		ptr := input.OutputPointer
		// TODO: Check if tx or output doesn't exist
		outputTx := transactions[ptr.TransactionHash]
		address := outputTx.Outputs[ptr.OutputIndex].ReceiverAddress

		// Find matching utxo
		utxoInd := -1
		for j, utxoPtr := range unspentOutputs[address] {
			if ptr.Equal(utxoPtr) {
				utxoInd = j
			}
		}

		// No utxo pointer matches
		if utxoInd == -1 {
			return false
		}

		// Remove utxo pointer at index utxoInd
		unspentOutputs[address][utxoInd] = unspentOutputs[address][len(unspentOutputs[address])-1]
		unspentOutputs[address] = unspentOutputs[address][:len(unspentOutputs[address])-1]
	}

	for outputIndex, output := range tx.Outputs {
		// Put new utxo pointer in
		outputPointer := transaction.TransactionOutputPointer{
			TransactionHash: txHash,
			OutputIndex:     uint16(outputIndex),
		}
		receiverAddress := output.ReceiverAddress
		unspentOutputs[receiverAddress] = append(unspentOutputs[receiverAddress], outputPointer)
	}

	return true
}

// Add a block to the chain
// TODO: Make this atomic - either it updates entire state if success or not at all
// Assumes that block has been validated by node
func AddBlock(block block.Block) (success bool) {
	// Add block to blocks
	// Update block hash
	// Process all transactions
	// Update pending transactions
	// Process all new and deleted utxos

	for _, tx := range block.Body {
		success := AddTransaction(tx)
		// Make sure transactions were successfully added
		if !success {
			return false
		}
	}

	// Add new block
	hash := block.Hash()
	blocks[hash] = block

	// Update last block hash
	lastBlockHash = hash

	return true
}

// Get all unspent output pointers for a given address
func GetUnspentTransactions(address [protocol.AddressLength]byte) (outputPointers []transaction.TransactionOutputPointer, success bool) {
	return unspentOutputs[address], true
}

// Get the output amount corresponding to a specific output pointer
func GetOutputAmount(ptr transaction.TransactionOutputPointer) (amount uint64, success bool) {
	hash := ptr.TransactionHash
	index := ptr.OutputIndex

	tx := transactions[hash]
	output := tx.Outputs[index]

	return output.Amount, true
}
