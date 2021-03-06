package chain

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/block"
	"github.com/AndrewCLu/TestcoinNode/common"
	"github.com/AndrewCLu/TestcoinNode/transaction"
	"github.com/AndrewCLu/TestcoinNode/util"
)

type Chain struct {
	Blocks              map[common.Hash]*block.Block
	Transactions        map[common.Hash]*transaction.Transaction
	LastBlockHash       common.Hash
	PendingTransactions []*transaction.Transaction
	UnspentOutputs      map[common.Address][]*transaction.TransactionOutputPointer
}

// Sets up the initial state of the chain
func New() (chn *Chain, ok bool) {
	chain := Chain{
		Blocks:              make(map[common.Hash]*block.Block),
		Transactions:        make(map[common.Hash]*transaction.Transaction),
		LastBlockHash:       *new(common.Hash),
		PendingTransactions: []*transaction.Transaction{},
		UnspentOutputs:      make(map[common.Address][]*transaction.TransactionOutputPointer),
	}

	return &chain, true
}

// Initializes the chain with a genesis block, return a bool indicating success
func (chain *Chain) Initialize(genesisBlock *block.Block) bool {
	genesisBlockHash := genesisBlock.Hash()
	chain.Blocks[genesisBlockHash] = genesisBlock
	chain.LastBlockHash = genesisBlockHash

	chain.AddBlock(genesisBlock)

	return true
}

// Given a transaction hash, returns a pointer to the transaction
// Returns bool indicating success
func (chain *Chain) GetTransaction(hash common.Hash) (tx *transaction.Transaction, ok bool) {
	return chain.Transactions[hash], true
}

// Given a pending transaction, return its transaction fee
func (chain *Chain) GetPendingTransactionFee(tx *transaction.Transaction) (fee uint64, ok bool) {
	var inputs uint64 = 0
	var outputs uint64 = 0
	for _, input := range tx.Inputs {
		inputAmount, inputOk := chain.GetOutputAmount(input.OutputPointer)
		if !inputOk {
			return 0, false
		}
		inputs += inputAmount
	}
	for _, output := range tx.Outputs {
		outputs += output.Amount
	}

	return inputs - outputs, true
}

// Gets up to num pending transactions
// Returns bool indicating success
// TODO: Have a ranking of pending transactions to retrieve by time added or miner fee
func (chain *Chain) GetPendingTransactions(num int) (txs []*transaction.Transaction, ok bool) {
	txs = chain.PendingTransactions
	if num > len(txs) {
		return txs, true
	}

	return txs, true
}

// Gets up to num pending transactions
// Returns bool indicating success
func (chain *Chain) GetPendingTransactionsByAddress(address common.Address) (txs []*transaction.Transaction, ok bool) {
	for _, tx := range chain.PendingTransactions {
		for _, input := range tx.Inputs {
			outputTx := chain.Transactions[input.OutputPointer.TransactionHash]
			inputAddress := outputTx.Outputs[input.OutputPointer.OutputIndex].ReceiverAddress
			if inputAddress.Equal(address) {
				txs = append(txs, tx)
				break
			}
		}
	}

	return txs, true
}

// Add a pending transaction to the list
// Returns bool indicating success
func (chain *Chain) AddPendingTransaction(tx *transaction.Transaction) bool {
	chain.PendingTransactions = append(chain.PendingTransactions, tx)

	return true
}

// Removes a list of pending transactions from the pool
// Returns a bool indicating success
func (chain *Chain) RemovePendingTransactions(txs []*transaction.Transaction) bool {
	for _, tx := range txs {
		ind := -1
		for i, ptx := range chain.PendingTransactions {
			if tx.Equal(ptx) {
				ind = i
			}
		}

		if ind != -1 {
			// Remove pending transaction at index ind
			chain.PendingTransactions[ind] = chain.PendingTransactions[len(chain.PendingTransactions)-1]
			chain.PendingTransactions = chain.PendingTransactions[:len(chain.PendingTransactions)-1]
		}
	}

	return true
}

// Get information about the last block in the chain
// Gets the hash and block number of the last blcok
// Returns bool indicating success
func (chain *Chain) GetLastBlockInfo() (hash common.Hash, blockNum int, ok bool) {
	return chain.LastBlockHash, len(chain.Blocks) - 1, true
}

// Adds a confirmed transaction to the chain
// Returns bool indicating success
// This is not a smart function - it will add the transaction, update the pending transactions and the utxos without validation
func (chain *Chain) AddTransaction(tx *transaction.Transaction) (ok bool) {
	// Find matching pending transaction
	ind := -1
	for i, ptx := range chain.PendingTransactions {
		if tx.Equal(ptx) {
			ind = i
		}
	}

	if ind != -1 {
		// Remove pending transaction at index ind
		chain.PendingTransactions[ind] = chain.PendingTransactions[len(chain.PendingTransactions)-1]
		chain.PendingTransactions = chain.PendingTransactions[:len(chain.PendingTransactions)-1]
	}

	// Add new transaction
	txHash := tx.Hash()
	chain.Transactions[txHash] = tx

	for _, input := range tx.Inputs {
		ptr := input.OutputPointer
		outputTx := chain.Transactions[ptr.TransactionHash]
		address := outputTx.Outputs[ptr.OutputIndex].ReceiverAddress

		// Find matching utxo
		utxoInd := -1
		for j, utxoPtr := range chain.UnspentOutputs[address] {
			if ptr.Equal(utxoPtr) {
				utxoInd = j
			}
		}

		if utxoInd != -1 {
			// Remove utxo pointer at index utxoInd
			chain.UnspentOutputs[address][utxoInd] = chain.UnspentOutputs[address][len(chain.UnspentOutputs[address])-1]
			chain.UnspentOutputs[address] = chain.UnspentOutputs[address][:len(chain.UnspentOutputs[address])-1]
		}
	}

	for outputIndex, output := range tx.Outputs {
		// Put new utxo pointer in
		outputPointer := transaction.TransactionOutputPointer{
			TransactionHash: txHash,
			OutputIndex:     uint16(outputIndex),
		}
		receiverAddress := output.ReceiverAddress
		chain.UnspentOutputs[receiverAddress] = append(chain.UnspentOutputs[receiverAddress], &outputPointer)
	}

	return true
}

// Add a block to the chain
// Returns bool indicating success
// This is not a smart function - it will add the block without validation
func (chain *Chain) AddBlock(block *block.Block) (ok bool) {
	for _, tx := range block.Body {
		fmt.Printf("Adding transaction to chain %v\n", tx.Hash().Hex())
		chain.AddTransaction(tx)
	}
	fmt.Printf("Adding coinbase to chain: %v\n", block.Coinbase.Hash().Hex())
	chain.AddTransaction(block.Coinbase)

	// Add new block
	hash := block.Hash()
	chain.Blocks[hash] = block

	// Update last block hash
	chain.LastBlockHash = hash

	return true
}

// Get all unspent output pointers for a given address
// Returns bool indicating success
func (chain *Chain) GetUnspentTransactions(address common.Address) (outputPointers []*transaction.TransactionOutputPointer, ok bool) {
	return chain.UnspentOutputs[address], true
}

// Get the output amount corresponding to a specific output pointer
// Returns bool indicating success
func (chain *Chain) GetOutputAmount(ptr *transaction.TransactionOutputPointer) (amount uint64, success bool) {
	hash := ptr.TransactionHash
	index := ptr.OutputIndex

	tx := chain.Transactions[hash]
	output := tx.Outputs[index]

	return output.Amount, true
}

// Gets the value of an account based on an address
func (chain *Chain) GetAccountValue(address common.Address) uint64 {
	var total uint64 = 0
	outputPointers, _ := chain.GetUnspentTransactions(address)
	for _, ptr := range outputPointers {
		amount, _ := chain.GetOutputAmount(ptr)
		total += amount
	}

	return total
}

// Prints the current state of the blockchain
func (chain *Chain) PrintChainState() {
	fmt.Printf("-------------------PRINTING CHAIN STATE-------------------\n")
	fmt.Printf("Blocks mined...\n")
	for _, block := range chain.Blocks {
		fmt.Printf("Block: %v\n", block.Hash().Hex())
	}

	fmt.Printf("Transactions confirmed...\n")
	for _, tx := range chain.Transactions {
		fmt.Printf("Transaction: %v\n", tx.Hash().Hex())
	}

	fmt.Printf("Unspent transactions...\n")
	for address, outputList := range chain.UnspentOutputs {
		amount := util.Uint64UnitToFloat64Unit(chain.GetAccountValue(address))
		fmt.Printf("Account %v has value %v\n", address.Hex(), amount)
		for _, output := range outputList {
			fmt.Printf("Account %v has unspent output at transaction %v index %v\n",
				address.Hex(),
				output.TransactionHash.Hex(),
				output.OutputIndex,
			)
		}
	}

	fmt.Printf("There are %v pending trnasactions\n", len(chain.PendingTransactions))
	fmt.Printf("-------------------END CHAIN STATE-------------------\n")
}

// Shallow copies chain over into another chain
func (chain *Chain) UnsafeCopy() *Chain {
	otherChain, _ := New()
	for blockHash, block := range chain.Blocks {
		otherChain.Blocks[blockHash] = block
	}
	for txHash, tx := range chain.Transactions {
		otherChain.Transactions[txHash] = tx
	}
	otherChain.LastBlockHash = chain.LastBlockHash
	for _, tx := range chain.PendingTransactions {
		otherChain.PendingTransactions = append(otherChain.PendingTransactions, tx)
	}
	for address, utxos := range chain.UnspentOutputs {
		newUtxos := []*transaction.TransactionOutputPointer{}
		for _, utxo := range utxos {
			newUtxos = append(newUtxos, utxo)
		}
		otherChain.UnspentOutputs[address] = newUtxos
	}

	return otherChain
}
