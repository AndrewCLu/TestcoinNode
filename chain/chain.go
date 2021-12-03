package chain

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/block"
	"github.com/AndrewCLu/TestcoinNode/common"
	"github.com/AndrewCLu/TestcoinNode/transaction"
	"github.com/AndrewCLu/TestcoinNode/util"
)

var blocks map[common.Hash]*block.Block                   // Stores all previous blocks
var transactions map[common.Hash]*transaction.Transaction // Stores all previous transactions

var lastBlockHash common.Hash                                                 // Last block hash
var pendingTransactions []*transaction.Transaction                            // All transactions that have not been processed yet
var unspentOutputs map[common.Address][]*transaction.TransactionOutputPointer // All utxos

// Sets up the initial state of the chain
func InitializeChain() {
	// Set up the genesis block
	genesisBlock := block.GetGenesisBlock()
	genesisBlockHash := genesisBlock.Hash()
	blocks = make(map[common.Hash]*block.Block)
	blocks[genesisBlockHash] = genesisBlock
	lastBlockHash = genesisBlockHash

	pendingTransactions = []*transaction.Transaction{}
	transactions = make(map[common.Hash]*transaction.Transaction)
	unspentOutputs = make(map[common.Address][]*transaction.TransactionOutputPointer)
}

// Given a transaction hash, returns a pointer to the transaction
// Returns bool indicating success
func GetTransaction(hash common.Hash) (tx *transaction.Transaction, ok bool) {
	return transactions[hash], true
}

// Gets up to num pending transactions
// Returns bool indicating success
// TODO: Have a ranking of pending transactions to retrieve by time added or miner fee
func GetPendingTransactions(num int) (txs []*transaction.Transaction, ok bool) {
	if num > len(pendingTransactions) {
		return pendingTransactions, true
	}

	return pendingTransactions[:num], true
}

// Add a pending transaction to the list
// Returns bool indicating success
func AddPendingTransaction(tx *transaction.Transaction) (ok bool) {
	pendingTransactions = append(pendingTransactions, tx)

	return true
}

// Get information about the last block in the chain
// Gets the hash and block number of the last blcok
// Returns bool indicating success
func GetLastBlockInfo() (hash common.Hash, blockNum int, ok bool) {
	return lastBlockHash, len(blocks) - 1, true
}

// Adds a confirmed transaction to the chain
// Returns bool indicating success
// TODO: Make this atomic - either it updates entire state if success or not at all
func AddTransaction(tx *transaction.Transaction) (ok bool) {
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
		unspentOutputs[receiverAddress] = append(unspentOutputs[receiverAddress], &outputPointer)
	}

	return true
}

// Add a block to the chain
// Assumes that block has been validated by node
// Returns bool indicating success
// TODO: Make this atomic - either it updates entire state if success or not at all
func AddBlock(block *block.Block) (ok bool) {
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
// Returns bool indicating success
func GetUnspentTransactions(address common.Address) (outputPointers []*transaction.TransactionOutputPointer, ok bool) {
	return unspentOutputs[address], true
}

// Get the output amount corresponding to a specific output pointer
// Returns bool indicating success
func GetOutputAmount(ptr *transaction.TransactionOutputPointer) (amount uint64, success bool) {
	hash := ptr.TransactionHash
	index := ptr.OutputIndex

	tx := transactions[hash]
	output := tx.Outputs[index]

	return output.Amount, true
}

// Gets the value of an account based on an address
func GetAccountValue(address common.Address) uint64 {
	var total uint64 = 0
	outputPointers, _ := GetUnspentTransactions(address)
	for _, ptr := range outputPointers {
		amount, _ := GetOutputAmount(ptr)
		total += amount
	}

	return total
}

// Prints the current state of the blockchain
func PrintChainState() {
	fmt.Printf("Blocks mined...\n")
	for _, block := range blocks {
		fmt.Printf("Block: %v\n", block.Hash().Hex())
	}

	fmt.Printf("Transactions confirmed...\n")
	for _, tx := range transactions {
		fmt.Printf("Transaction: %v\n", tx.Hash().Hex())
	}

	fmt.Printf("Unspent transactions...\n")
	for address, outputList := range unspentOutputs {
		amount := util.Uint64UnitToFloat64Unit(GetAccountValue(address))
		fmt.Printf("Account %v has value %v\n", address.Hex(), amount)
		for _, output := range outputList {
			fmt.Printf("Account %v has unspent output at transaction %v index %v\n",
				address.Hex(),
				output.TransactionHash.Hex(),
				output.OutputIndex,
			)
		}
	}
}
