package node

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/block"
	"github.com/AndrewCLu/TestcoinNode/chain"
	"github.com/AndrewCLu/TestcoinNode/common"
	"github.com/AndrewCLu/TestcoinNode/consensus"
	"github.com/AndrewCLu/TestcoinNode/consensus/pow"
	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/miner"
	"github.com/AndrewCLu/TestcoinNode/protocol"
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

// Initializes the node by beginning the chain with the genesis block
func (node *Node) Initialize(coinbaseAddress common.Address) bool {
	genesisBlock := GetGenesisBlock(coinbaseAddress)
	chainOk := node.Chain.Initialize(genesisBlock)

	return chainOk
}

// Returns a new account
// TODO: Key management for accounts
func (node *Node) NewAccount() *account.Account {
	account, _ := account.New()
	fmt.Printf("Created account with address: %v\n", account.Address.Hex())

	return account
}

// Returns a pointer to hard coded genesis block
func GetGenesisBlock(coinbaseAddress common.Address) *block.Block {
	coinbaseOutput := &transaction.TransactionOutput{
		ReceiverAddress: coinbaseAddress,
		Amount:          protocol.ComputeBlockReward(0),
	}
	coinbase, _ := transaction.New(
		[]*transaction.TransactionInput{},
		[]*transaction.TransactionOutput{coinbaseOutput},
	)
	block, _ := block.New(crypto.HashBytes([]byte("first")), 0, []*transaction.Transaction{}, coinbase)

	return block
}

// Validates a transaction and if valid, adds it to the chain's pool of pending transactions
// Returns a bool indicating success
func (node *Node) AddPendingTransaction(tx *transaction.Transaction) bool {
	validateTx := node.Consensus.ValidatePendingTransaction(node.Chain, tx)
	if !validateTx {
		fmt.Println("Failed to validate new transaction, not adding to chain")
		return false
	}

	node.Chain.AddPendingTransaction(tx)
	return true
}

// Creates a new coinbase transaction for a given account
// Testing function only, this is only ever created by the miner
func (node *Node) NewCoinbaseTransaction(account *account.Account, readableAmount float64) *transaction.Transaction {
	address := account.Address
	amount := util.Float64UnitToUnit64Unit(readableAmount)

	output := &transaction.TransactionOutput{ReceiverAddress: address, Amount: amount}
	newTransaction, success := transaction.New(
		[]*transaction.TransactionInput{},
		[]*transaction.TransactionOutput{output},
	)

	if !success || !node.Consensus.ValidatePendingTransaction(node.Chain, newTransaction) {
		fmt.Println("Attempted to create new coinbase transaction and FAILED")
		return nil
	}

	fmt.Printf("Created new coinbase transaction %v sending %v to %v\n",
		newTransaction.Hash().Hex(),
		readableAmount,
		address.Hex(),
	)

	return newTransaction
}

// Creates a new peer transaction for a given amount
// Readable indicates that the units taken in by this function are in decimal units, which need to be converted to integer units before sending
func (node *Node) NewPeerTransaction(account *account.Account, receiverAddress common.Address, readableAmount float64, readableTransactionFee float64) *transaction.Transaction {
	senderAddress := account.Address
	senderPublicKey := account.PublicKey
	senderPrivateKey := account.PrivateKey
	amount := util.Float64UnitToUnit64Unit(readableAmount)
	transactionFee := util.Float64UnitToUnit64Unit(readableTransactionFee)

	// Check that sender has enough money
	senderValue := node.Chain.GetAccountValue(senderAddress)
	if senderValue < amount+transactionFee {
		fmt.Println("Attempted to create new peer transaction but sender has insufficient funds.")
		return nil
	}

	allUtxos, _ := node.Chain.GetUnspentTransactions(senderAddress)
	pendingTransactions, _ := node.Chain.GetPendingTransactionsByAddress(senderAddress)
	pendingUtxos := []*transaction.TransactionOutputPointer{}
	selectedUtxos := []*transaction.TransactionOutputPointer{}
	var currentAmount uint64 = 0
	// Two step process: First select Utxos that are not used in a pending transaction
	// If all utxos are used, then select ones that are part of a pending transaction
finishedSelection:
	for _, utxo := range allUtxos {
		match := false
	nextUtxo:
		for _, tx := range pendingTransactions {
			for _, input := range tx.Inputs {
				// Utxo is already pending, don't add to this transaction unless we cannot achieve the desired amount otherwise
				if input.OutputPointer.Equal(utxo) {
					match = true
					pendingUtxos = append(pendingUtxos, utxo)
					break nextUtxo
				}
			}
		}
		// Utxo is not pending, add to this transaction
		if !match {
			selectedUtxos = append(selectedUtxos, utxo)
			newAmount, _ := node.Chain.GetOutputAmount(utxo)
			currentAmount += newAmount
		}

		if currentAmount >= amount+transactionFee {
			break finishedSelection
		}
	}
	// Add utxos that are part of a pending trnasaction if the amount has not been reached
	if currentAmount < amount+transactionFee {
		for _, utxo := range pendingUtxos {
			selectedUtxos = append(selectedUtxos, utxo)
			newAmount, _ := node.Chain.GetOutputAmount(utxo)
			currentAmount += newAmount

			if currentAmount >= amount+transactionFee {
				break
			}
		}
	}

	inputs := []*transaction.TransactionInput{}
	for _, ptr := range selectedUtxos {
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
	diff := currentAmount - amount - transactionFee

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

	if !success || !node.Consensus.ValidatePendingTransaction(node.Chain, newTransaction) {
		node.PrintTransaction(newTransaction)
		fmt.Println("Attempted to create new peer transaction and FAILED")
		return nil
	}

	fmt.Printf("Created new peer transaction %v sending %v from %v to %v with transaction fee of %v\n",
		newTransaction.Hash().Hex(),
		readableAmount,
		senderAddress.Hex(),
		receiverAddress.Hex(),
		readableTransactionFee,
	)

	node.AddPendingTransaction(newTransaction)
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
	block, ok := node.Miner.MineBlock()
	if !ok {
		return
	}

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

// Prints a given transaction
func (node *Node) PrintTransaction(tx *transaction.Transaction) {
	fmt.Printf("Printing transaction %v...", tx.Hash().Hex())
	fmt.Printf("Inputs: ")
	for _, input := range tx.Inputs {
		inputAmount, _ := node.Chain.GetOutputAmount(input.OutputPointer)
		fmt.Printf("%v from transaction %v ", inputAmount, input.OutputPointer.TransactionHash.Hex())
	}
	fmt.Printf("Outputs: ")
	for _, output := range tx.Outputs {
		fmt.Printf("%v to address %v ", output.Amount, output.ReceiverAddress.Hex())
	}
	fmt.Printf("\n")
}
