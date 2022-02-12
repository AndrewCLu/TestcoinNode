package node

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/chain"
	"github.com/AndrewCLu/TestcoinNode/common"
	"github.com/AndrewCLu/TestcoinNode/consensus"
	"github.com/AndrewCLu/TestcoinNode/consensus/pow"
	"github.com/AndrewCLu/TestcoinNode/miner"
	"github.com/AndrewCLu/TestcoinNode/transaction"
	"github.com/AndrewCLu/TestcoinNode/util"
)

type Node struct {
	Chain     *chain.Chain
	Consensus consensus.Consensus
	Miner     *miner.Miner
}

func New() (n *Node, ok bool) {
	chn, _ := chain.New()
	pow, _ := pow.New()
	node := Node{
		Chain:     chn,
		Consensus: pow,
	}

	return &node, true
}

// Returns a new account
// TODO: Key management for accounts
func (node *Node) NewAccount() *account.Account {
	account, _ := account.New()
	fmt.Printf("Created account with address: %v\n", account.Address.Hex())

	return account
}

// Creates a new coinbase transaction for a given account
func (node *Node) NewCoinbaseTransaction(account *account.Account, readableAmount float64) *transaction.Transaction {
	address := account.Address
	amount := util.Float64UnitToUnit64Unit(readableAmount)

	output := &transaction.TransactionOutput{ReceiverAddress: address, Amount: amount}
	newTransaction, success := transaction.New(
		[]*transaction.TransactionInput{},
		[]*transaction.TransactionOutput{output},
	)

	if !success || !node.Consensus.ValidateTransaction(node.Chain, newTransaction) {
		fmt.Printf("Attempted to create new coinbase transaction and FAILED")
		return nil
	}

	fmt.Printf("Created new coinbase transaction %v sending %v to %v\n",
		newTransaction.Hash().Hex(),
		readableAmount,
		address.Hex(),
	)
	node.Chain.AddPendingTransaction(newTransaction)
	return newTransaction
}

// Creates a new peer transaction for a given amount
func (node *Node) NewPeerTransaction(account *account.Account, receiverAddress common.Address, readableAmount float64) *transaction.Transaction {
	senderAddress := account.Address
	senderPublicKey := account.PublicKey
	senderPrivateKey := account.PrivateKey
	amount := util.Float64UnitToUnit64Unit(readableAmount)

	// Check that sender has enough money
	senderValue := node.Chain.GetAccountValue(senderAddress)
	if senderValue < amount {
		fmt.Printf("Attempted to create new peer transaction but sender has insufficient funds.")
		return nil
	}

	outputPointers, _ := node.Chain.GetUnspentTransactions(senderAddress)

	// Current implementation just uses all utxos in a transaction
	// TODO: Pick the minimum number of utxos an account can use to complete a transaction
	inputs := []*transaction.TransactionInput{}
	for _, ptr := range outputPointers {
		signature := node.Consensus.SignInput(senderPrivateKey, ptr)

		verification := &transaction.TransactionInputVerification{
			Signature:        signature,
			EncodedPublicKey: senderPublicKey,
		}

		input := &transaction.TransactionInput{
			OutputPointer:      ptr,
			VerificationLength: uint16(len(verification.Bytes())),
			Verification:       verification,
		}

		inputs = append(inputs, input)
	}

	// If sender has more money than amount, create a refund transaction output
	diff := senderValue - amount

	outputReceiver := &transaction.TransactionOutput{
		ReceiverAddress: receiverAddress,
		Amount:          amount,
	}
	outputSender := &transaction.TransactionOutput{
		ReceiverAddress: senderAddress,
		Amount:          diff,
	}

	outputs := []*transaction.TransactionOutput{outputReceiver}
	if diff != 0 {
		outputs = append(outputs, outputSender)
	}

	newTransaction, success := transaction.New(inputs, outputs)

	if !success || !node.Consensus.ValidateTransaction(node.Chain, newTransaction) {
		fmt.Printf("Attempted to create new peer transaction and FAILED")
		return nil
	}

	fmt.Printf("Created new peer transaction %v sending %v from %v to %v\n",
		newTransaction.Hash().Hex(),
		readableAmount,
		senderAddress.Hex(),
		receiverAddress.Hex(),
	)
	node.Chain.AddPendingTransaction(newTransaction)
	return newTransaction
}

// Initializes the miner with specified coinbase address
func (node *Node) BeginMiner(coinbase common.Address) {
	node.Miner, _ = miner.New(coinbase, node.Chain, node.Consensus)
}

// Calls the miner to mine a block and adds it to the chain if it is valid
func (node *Node) MineBlock() {
	if node.Miner == nil {
		return
	}
	block, _ := node.Miner.MineBlock()
	valid := node.Consensus.ValidateBlock(node.Chain, block)

	if valid {
		node.Chain.AddBlock(block)
	}
}

// Gets the human readable value of an account
func (node *Node) GetReadableAccountValue(account *account.Account) float64 {
	address := account.Address

	total := node.Chain.GetAccountValue(address)
	readableTotal := util.Uint64UnitToFloat64Unit(total)

	fmt.Printf("Account with address %v has value %v\n", address.Hex(), readableTotal)

	return readableTotal
}

// Testing function to print the state of the chain
func (node *Node) PrintChainState() {
	node.Chain.PrintChainState()
}
