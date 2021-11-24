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
