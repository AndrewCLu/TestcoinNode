package node

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/crypto"
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

	if !success || !ValidateTransaction(account.GetPublicKey(), newTransaction) {
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

	newTransaction, success := transaction.NewPeerTransaction(account.GetPrivateKey(), utxos, outputs)

	if !success || !ValidateTransaction(account.GetPublicKey(), newTransaction) {
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

// Returns if a transaction is valid or not based on the state of the ledger
// TODO: Make verifications of utxos existing in ledger more efficient
// TODO: Safety checks
func ValidateTransaction(senderPublicKey []byte, tx transaction.Transaction) bool {
	senderAddress := crypto.GetAddressFromPublicKey(senderPublicKey)

	var inputTotal uint64 = 0
	for _, input := range tx.Inputs {
		hash := input.PreviousTransactionHash
		index := input.PreviousTransactionIndex
		signature := input.SenderSignature

		previousTransaction := ledger[hash]
		previousTransactionOutput := previousTransaction.Outputs[index]
		utxo := transaction.UnspentTransactionOutput{
			TransactionHash:  hash,
			TransactionIndex: index,
			ReceiverAddress:  previousTransactionOutput.ReceiverAddress,
			Amount:           previousTransactionOutput.Amount,
		}

		// Check if input is provided by the sender
		if !transaction.VerifyInput(senderPublicKey, hash, index, signature) {
			return false
		}

		// Find if a current utxo matches the one implied by the transaction
		match := false
		for _, compareUtxo := range unspentOutputs[senderAddress] {
			if utxo.Equal(compareUtxo) {
				match = true
			}
		}

		if !match {
			return false
		}

		inputTotal += utxo.Amount
	}

	var outputTotal uint64 = 0
	for _, output := range tx.Outputs {
		outputTotal += output.Amount
	}

	if len(tx.Inputs) != 0 && inputTotal != outputTotal {
		return false
	}

	return true
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
