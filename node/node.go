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
var unspentOutputs map[[protocol.AddressLength]byte][]transaction.UnspentTransactionOutput

func InitializeNode() {
	ledger = make(map[[transaction.TransactionHashLength]byte]transaction.Transaction)
	unspentOutputs = make(map[[protocol.AddressLength]byte][]transaction.UnspentTransactionOutput)
}

// Creates a new coinbase transaction for a given account
func NewCoinbaseTransaction(account account.Account, readableAmount float64) {
	address := account.GetAddress()
	amount := util.Float64UnitToUnit64Unit(readableAmount)
	newTransaction, success := transaction.NewCoinbaseTransaction(address, amount)
	transactionHash := newTransaction.GetTransactionHash()

	if !success {
		return
	}

	// Add transaction to ledger
	ledger[transactionHash] = newTransaction

	// Record unspent transaction output
	output := transaction.UnspentTransactionOutput{
		TransactionHash:  transactionHash,
		TransactionIndex: uint16(0),
		ReceiverAddress:  address,
		Amount:           amount,
	}
	unspentOutputs[address] = append(unspentOutputs[address], output)

	fmt.Printf("Coinbase transaction %v sending %v to %v\n", transactionHash, readableAmount, address)
}

// Creates a new peer transaction for a given amount
func NewPeerTransaction(account account.Account, receiverAddress [protocol.AddressLength]byte, readableAmount float64) {
	senderAddress := account.GetAddress()
	amount := util.Float64UnitToUnit64Unit(readableAmount)

	utxos := unspentOutputs[senderAddress]
	output := transaction.TransactionOutput{
		ReceiverAddress: receiverAddress,
		Amount:          amount,
	}
	outputs := []transaction.TransactionOutput{output}

	newTransaction, success := transaction.NewPeerTransaction(account.GetPrivateKey(), utxos, outputs)
	transactionHash := newTransaction.GetTransactionHash()

	if !success {
		return
	}

	// Add transaction to ledger
	ledger[transactionHash] = newTransaction

	// Record unspent transaction output
	newOutput := transaction.UnspentTransactionOutput{
		TransactionHash:  transactionHash,
		TransactionIndex: uint16(0),
		ReceiverAddress:  receiverAddress,
		Amount:           amount,
	}
	unspentOutputs[receiverAddress] = append(unspentOutputs[receiverAddress], newOutput)

	fmt.Printf("Coinbase transaction %v sending %v from %v to %v\n", transactionHash, readableAmount, senderAddress, receiverAddress)
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
