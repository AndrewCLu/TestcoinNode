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

func GetTransaction(hash [crypto.HashLength]byte) (tx transaction.Transaction, success bool) {
	return transactions[hash], true
}

func GetPendingTransactions(num int) (txs []transaction.Transaction, success bool) {
	if num > len(pendingTransactions) {
		return pendingTransactions, true
	}

	return pendingTransactions[:num], true
}

func AddPendingTransaction(tx transaction.Transaction) (success bool) {
	pendingTransactions = append(pendingTransactions, tx)

	return true
}

func GetLastBlockInfo() (hash [crypto.HashLength]byte, blockNum int, success bool) {
	return lastBlockHash, len(blocks) - 1, true
}

func AddBlock(block block.Block) (success bool) {
	// Add block to blocks
	// Update block hash
	// Process all transactions
	// Update pending transactions
	// Process all new and deleted utxos

	return true
}

func GetUnspentTransactions(address [protocol.AddressLength]byte) (outputPointers []transaction.TransactionOutputPointer, success bool) {
	return unspentOutputs[address], true
}

func GetOutputAmount(ptr transaction.TransactionOutputPointer) (amount uint64, success bool) {
	hash := ptr.TransactionHash
	index := ptr.OutputIndex

	tx := transactions[hash]
	output := tx.Outputs[index]

	return output.Amount, true
}

// ledger[transactionHash] = newTransaction

// 	// Record unspent transaction output
// 	output := transaction.UnspentTransactionOutput{
// 		TransactionHash:  transactionHash,
// 		TransactionIndex: uint16(0),
// 		ReceiverAddress:  address,
// 		Amount:           amount,
// 	}
// 	unspentOutputs[address] = append(unspentOutputs[address], output)

// fmt.Printf("Coinbase transaction %v sending %v to %v\n",
// 	util.HashToHexString(transactionHash),
// 	readableAmount,
// 	util.AddressToHexString(address),

// // Add transaction to ledger
// 	ledger[transactionHash] = newTransaction

// 	// Record received unspent transaction output
// 	receiverOutput := transaction.UnspentTransactionOutput{
// 		TransactionHash:  transactionHash,
// 		TransactionIndex: uint16(0),
// 		ReceiverAddress:  receiverAddress,
// 		Amount:           amount,
// 	}
// 	unspentOutputs[receiverAddress] = append(unspentOutputs[receiverAddress], receiverOutput)

// 	// Record refund to sender as unspent transaction output
// 	if diff != 0 {
// 		senderOutput := transaction.UnspentTransactionOutput{
// 			TransactionHash:  transactionHash,
// 			TransactionIndex: uint16(1),
// 			ReceiverAddress:  senderAddress,
// 			Amount:           diff,
// 		}
// 		unspentOutputs[senderAddress] = []transaction.UnspentTransactionOutput{senderOutput}
// 	}

// 	fmt.Printf("Peer transaction %v sending %v from %v to %v\n",
// 		util.HashToHexString(transactionHash),
// 		readableAmount,
// 		util.AddressToHexString(senderAddress),
// 		util.AddressToHexString(receiverAddress),
// 	)
