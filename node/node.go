package node

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/chain"
	"github.com/AndrewCLu/TestcoinNode/miner"
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
func NewCoinbaseTransaction(account account.Account, readableAmount float64) transaction.Transaction {
	address := account.GetAddress()
	amount := util.Float64UnitToUnit64Unit(readableAmount)

	output := transaction.TransactionOutput{ReceiverAddress: address, Amount: amount}
	newTransaction, success := transaction.NewTransaction(
		[]transaction.TransactionInput{},
		[]transaction.TransactionOutput{output},
	)

	if !success || !ValidateTransaction(newTransaction) {
		fmt.Printf("Attempted to create new coinbase transaction and FAILED")
		return transaction.Transaction{}
	}

	fmt.Printf("Created new coinbase transaction sending %v to %v\n", readableAmount, util.HashToHexString(address))
	chain.AddPendingTransaction(newTransaction)
	return newTransaction
}

// Creates a new peer transaction for a given amount
func NewPeerTransaction(account account.Account, receiverAddress [protocol.AddressLength]byte, readableAmount float64) transaction.Transaction {
	senderAddress := account.GetAddress()
	senderPublicKey := account.GetPublicKey()
	senderPrivateKey := account.GetPrivateKey()
	amount := util.Float64UnitToUnit64Unit(readableAmount)

	// Check that sender has enough money
	senderValue := GetAccountValue(senderAddress)
	if senderValue < amount {
		fmt.Printf("Attempted to create new peer transaction but sender has insufficient funds.")
		return transaction.Transaction{}
	}

	outputPointers, _ := chain.GetUnspentTransactions(senderAddress)

	// Current implementation just uses all utxos in a transaction
	// TODO: Pick the minimum number of utxos an account can use to complete a transaction
	inputs := []transaction.TransactionInput{}
	for _, ptr := range outputPointers {
		signature := SignInput(senderPrivateKey, ptr)

		verification := transaction.TransactionInputVerification{
			Signature:        signature,
			EncodedPublicKey: senderPublicKey,
		}

		input := transaction.TransactionInput{
			OutputPointer:      ptr,
			VerificationLength: uint16(len(verification.TransactionInputVerificationToByteArray())),
			Verification:       verification,
		}

		inputs = append(inputs, input)
	}

	// If sender has more money than amount, create a refund transaction output
	diff := senderValue - amount

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

	newTransaction, success := transaction.NewTransaction(inputs, outputs)

	if !success || !ValidateTransaction(newTransaction) {
		fmt.Printf("Attempted to create new peer transaction and FAILED")
		return transaction.Transaction{}
	}

	fmt.Printf("Created new peer transaction sending %v from %v to %v", readableAmount, util.HashToHexString(senderAddress), util.HashToHexString(receiverAddress))
	chain.AddPendingTransaction(newTransaction)
	return newTransaction
}

// Calls the miner to mine a block and adds it to the chain if it is valid
func MineBlock() {
	block := miner.MineBlock()
	valid := ValidateBlock(block)

	if valid {
		chain.AddBlock(block)
	}
}

// Gets the human readable value of an account
func GetReadableAccountValue(account account.Account) float64 {
	address := account.GetAddress()

	total := chain.GetAccountValue(address)
	readableTotal := util.Uint64UnitToFloat64Unit(total)

	fmt.Printf("Account with address %v has value %v\n", util.AddressToHexString(address), readableTotal)

	return readableTotal
}

// Testing function to print the state of the chain
func PrintChainState() {
	chain.PrintChainState()
}
