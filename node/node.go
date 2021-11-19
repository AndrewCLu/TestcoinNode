package node

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/protocol"
	"github.com/AndrewCLu/TestcoinNode/transaction"
	"github.com/AndrewCLu/TestcoinNode/util"
)

// A ledger mapping transaction hashes to transactions
var ledger map[[transaction.TransactionHashLength]byte]transaction.Transaction
var unspentOutputs map[[protocol.AddressLength]byte][]transaction.TransactionOutput

func InitializeNode() {
	ledger = make(map[[transaction.TransactionHashLength]byte]transaction.Transaction)
	unspentOutputs = make(map[[protocol.AddressLength]byte][]transaction.TransactionOutput)
}

// Creates a new coinbase transaction for a given account
func NewCoinbaseTransaction(account account.Account, readableAmount float64) {
	address := account.GetAddress()
	amount := util.Float64UnitToUnit64Unit(readableAmount)
	transaction, _ := transaction.NewCoinbaseTransaction(address, amount)
	transactionHash := transaction.GetTransactionHash()

	// Add transaction to ledger
	ledger[transactionHash] = transaction

	// Record unspent transaction output
	unspentOutputs[address] = append(unspentOutputs[address], transaction.Outputs[0])

	fmt.Printf("Coinbase transaction %v sending %v to %v\n", transactionHash, readableAmount, address)
}

// Creates a new peer transaction for a given amount
func NewPeerTransaction(account account.Account, receiverAddress [protocol.AddressLength]byte, rawAmount float64) {
	return
}

func GetAccountValue(account account.Account) float64 {
	address := account.GetAddress()

	total := uint64(0)
	for _, output := range unspentOutputs[address] {
		total += output.Amount
	}

	readableTotal := util.Uint64UnitToFloat64Unit(total)

	fmt.Printf("Account with address %v has value %v\n", address, readableTotal)
	return readableTotal
}
