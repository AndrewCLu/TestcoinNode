package chain

import (
	"github.com/AndrewCLu/TestcoinNode/block"
	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/protocol"
	"github.com/AndrewCLu/TestcoinNode/transaction"
)

var blocks map[[crypto.HashLength]byte]block.Block                                   // Stores all previous blocks
var transactions map[[transaction.TransactionHashLength]byte]transaction.Transaction // Stores all previous transactions

var lastBlockHash [crypto.HashLength]byte                                                  // Last block hash
var pendingTransactions []transaction.Transaction                                          // All transactions that have not been processed yet
var unspentOutputs map[[protocol.AddressLength]byte][]transaction.UnspentTransactionOutput // All utxos

func InitializeChain() {
	genesisBlock, _ := block.NewBlock(0, crypto.HashBytes([]byte("first")), []transaction.Transaction{})
	genesisBlockHash := genesisBlock.Hash()
	blocks = make(map[[crypto.HashLength]byte]block.Block)
	blocks[genesisBlockHash] = genesisBlock
	lastBlockHash = genesisBlockHash

	pendingTransactions = []transaction.Transaction{}
	transactions = make(map[[transaction.TransactionHashLength]byte]transaction.Transaction)
	unspentOutputs = make(map[[protocol.AddressLength]byte][]transaction.UnspentTransactionOutput)
}

func AddTransaction(tx transaction.Transaction) {
	// TODO: Validate transaction
	pendingTransactions = append(pendingTransactions, tx)
}

func GetUnspentTransactions(address [protocol.AddressLength]byte) []transaction.UnspentTransactionOutput {
	return unspentOutputs[address]
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
