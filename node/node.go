package node

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/chain"
	"github.com/AndrewCLu/TestcoinNode/protocol"
	"github.com/AndrewCLu/TestcoinNode/transaction"
	"github.com/AndrewCLu/TestcoinNode/util"
)

func InitializeNode() {
	chain.InitializeChain()
}

// Returns a new account
func NewAccount() account.Account {
	account := account.NewAccount()

	fmt.Printf("Created account with address: %v\n", util.AddressToHexString(account.GetAddress()))

	return account
}

// Creates a new coinbase transaction for a given account
func NewCoinbaseTransaction(account account.Account, readableAmount float64) {
	address := account.GetAddress()
	amount := util.Float64UnitToUnit64Unit(readableAmount)
	newTransaction, success := transaction.NewCoinbaseTransaction(address, amount)

	if !success || !ValidateTransaction(newTransaction) {
		fmt.Printf("Attempted to create new coinbase transaction and FAILED")
		return
	}

	transactionHash := newTransaction.Hash()

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

	fmt.Printf("Coinbase transaction %v sending %v to %v\n",
		util.HashToHexString(transactionHash),
		readableAmount,
		util.AddressToHexString(address),
	)
}

// Creates a new peer transaction for a given amount
func NewPeerTransaction(account account.Account, receiverAddress [protocol.AddressLength]byte, readableAmount float64) {
	senderAddress := account.GetAddress()
	amount := util.Float64UnitToUnit64Unit(readableAmount)

	senderValue := GetAccountValue(senderAddress)
	if senderValue < amount {
		return
	}

	diff := senderValue - amount

	utxos := unspentOutputs[senderAddress]
	outputReceiver := transaction.TransactionOutput{
		ReceiverAddress: receiverAddress,
		Amount:          amount,
	}
	outputSender := transaction.TransactionOutput{
		ReceiverAddress: senderAddress,
		Amount:          diff,
	}

	outputs := []transaction.TransactionOutput{outputReceiver}
	if diff != 0 {
		outputs = append(outputs, outputSender)
	}

	newTransaction, success := transaction.NewPeerTransaction(account.GetPublicKey(), account.GetPrivateKey(), utxos, outputs)

	if !success || !ValidateTransaction(newTransaction) {
		fmt.Printf("Attempted to create new peer transaction and FAILED")
		return
	}

	transactionHash := newTransaction.Hash()

	// Add transaction to ledger
	ledger[transactionHash] = newTransaction

	// Record received unspent transaction output
	receiverOutput := transaction.UnspentTransactionOutput{
		TransactionHash:  transactionHash,
		TransactionIndex: uint16(0),
		ReceiverAddress:  receiverAddress,
		Amount:           amount,
	}
	unspentOutputs[receiverAddress] = append(unspentOutputs[receiverAddress], receiverOutput)

	// Record refund to sender as unspent transaction output
	if diff != 0 {
		senderOutput := transaction.UnspentTransactionOutput{
			TransactionHash:  transactionHash,
			TransactionIndex: uint16(1),
			ReceiverAddress:  senderAddress,
			Amount:           diff,
		}
		unspentOutputs[senderAddress] = []transaction.UnspentTransactionOutput{senderOutput}
	}

	fmt.Printf("Peer transaction %v sending %v from %v to %v\n",
		util.HashToHexString(transactionHash),
		readableAmount,
		util.AddressToHexString(senderAddress),
		util.AddressToHexString(receiverAddress),
	)
}

// Gets the micro unit value of an account based on an address
func GetAccountValue(address [protocol.AddressLength]byte) uint64 {
	var total uint64 = 0
	for _, output := range unspentOutputs[address] {
		total += output.Amount
	}

	return total
}

// Gets the human readable value of an account
func GetReadableAccountValue(account account.Account) float64 {
	address := account.GetAddress()

	total := GetAccountValue(address)
	readableTotal := util.Uint64UnitToFloat64Unit(total)

	fmt.Printf("Account with address %v has value %v\n", util.AddressToHexString(address), readableTotal)

	return readableTotal
}
