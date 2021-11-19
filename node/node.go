package node

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/protocol"
	"github.com/AndrewCLu/TestcoinNode/transaction"
)

// A ledger mapping transaction hashes to transactions
var ledger map[[transaction.TransactionHashLength]byte]transaction.Transaction
var unspentOutputs map[[protocol.AddressLength]byte][]transaction.TransactionOutput

func InitializeNode() {
	ledger = make(map[[transaction.TransactionHashLength]byte]transaction.Transaction)
	unspentOutputs = make(map[[protocol.AddressLength]byte][]transaction.TransactionOutput)
}

// Creates a new coinbase transaction for a given account
func NewCoinbaseTransaction(account account.Account, amount float64) {
	address := account.GetAddress()
	transaction := transaction.NewCoinbaseTransaction(address, amount)
	transactionHash := transaction.GetTransactionHash()

	// Add transaction to ledger
	ledger[transactionHash] = transaction

	// Record unspent transaction output
	unspentOutputs[address] = append(unspentOutputs[address], transaction.Outputs[0])

	fmt.Printf("Coinbase transaction %v sending %v to %v\n", transactionHash, amount, address)
}

// Creates a new peer transaction for a given amount
func NewPeerTransaction(account account.Account, receiverAddress [protocol.AddressLength]byte, amount float64) {
	return
}

func GetAccountValue(account account.Account) float64 {
	address := account.GetAddress()

	total := 0.0
	for _, output := range unspentOutputs[address] {
		total += output.Amount
	}

	fmt.Printf("Account with address %v has value %v\n", address, total)
	return total
}
