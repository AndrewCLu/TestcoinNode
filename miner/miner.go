package miner

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/AndrewCLu/TestcoinNode/block"
	"github.com/AndrewCLu/TestcoinNode/chain"
	"github.com/AndrewCLu/TestcoinNode/common"
	"github.com/AndrewCLu/TestcoinNode/consensus"
	"github.com/AndrewCLu/TestcoinNode/protocol"
	"github.com/AndrewCLu/TestcoinNode/transaction"
)

const DefaultHashLimit = 10 * 1000 * 1000 // The maximum number of hashes a miner will attempt to solve a block

// The configuration for a miner
type MinerConfig struct {
	HashLimit int
}

// Creates blocks and solves proof of work
type Miner struct {
	Coinbase  common.Address
	Config    *MinerConfig
	Chain     *chain.Chain
	Consensus consensus.Consensus
}

// Creates and returns the address of a new Miner
// Returns a bool indicating success
// TODO: Enable users to set up miners with non-default values
func New(coinbase common.Address, chain *chain.Chain, consensus consensus.Consensus) (*Miner, bool) {
	defaultConfig := MinerConfig{
		HashLimit: DefaultHashLimit,
	}

	miner := Miner{
		Coinbase:  coinbase,
		Config:    &defaultConfig,
		Chain:     chain,
		Consensus: consensus,
	}

	return &miner, true
}

// Tries to mine a block from pending transactions on the chain
// Returns the block and a boolean indicating success
// TODO: Miner must validate transactions
func (miner *Miner) MineBlock() (blk *block.Block, ok bool) {
	txs, txOk := miner.Chain.GetPendingTransactions(protocol.MaxTransactionsInBlock)
	if !txOk {
		fmt.Println("Could not get transactions from the current chain")
		return nil, false
	}

	// Sort transactions by decreasing fee
	sort.Slice(txs, func(i, j int) bool {
		iFee, _ := miner.Chain.GetPendingTransactionFee(txs[i])
		jFee, _ := miner.Chain.GetPendingTransactionFee(txs[j])
		return iFee > jFee
	})

	lastBlockHash, lastBlockNum, lastBlockOk := miner.Chain.GetLastBlockInfo()
	if !lastBlockOk {
		fmt.Println("Could not get last block from current chain")
		return nil, false
	}
	blockNum := lastBlockNum + 1

	coinbaseOutput := &transaction.TransactionOutput{
		ReceiverAddress: miner.Coinbase,
		Amount:          protocol.ComputeBlockReward(blockNum),
	}
	coinbase, _ := transaction.New(
		[]*transaction.TransactionInput{},
		[]*transaction.TransactionOutput{coinbaseOutput},
	)

	block, blockOk := block.New(lastBlockHash, blockNum, txs, coinbase)
	if !blockOk {
		fmt.Println("Failed to create new block")
		return nil, false
	}

	// Compute the nonce that solves a block
	nonce, solveOk := miner.solve(*block.Header)
	if !solveOk {
		fmt.Println("Failed to solve block with allotted parameters")
		return nil, false
	}
	block.Header.Nonce = nonce

	blockHash := block.Hash()
	for _, tx := range block.Body {
		fmt.Printf("Block %v confirmed transaction %v\n", blockHash.Hex(), tx.Hash().Hex())
	}

	return block, true
}

// Given a block header, compute the nonce that results in a hash under the desired target
// Does not modify the block header passed in
func (miner *Miner) solve(header block.BlockHeader) (nonce uint32, ok bool) {
	// Start the nonce at a random number to avoid multiple nodes mining the same nonces
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	nonce = r.Uint32()

	// Represent target as a full hash to compare block hash to
	target := header.Target.FullHash()

	t1 := time.Now()
	fmt.Printf("Solving block with target %v ...\n", target.Hex())

	// TODO: Check valid header using consensus
	for count := 0; count < miner.Config.HashLimit; count++ {
		count += 1
		header.Nonce = nonce
		hash := header.Hash()

		if bytes.Compare(hash.Bytes(), target.Bytes()) < 0 {
			t2 := time.Now()
			diff := t2.Sub(t1)

			fmt.Printf("Successfully found nonce %v with %v tries in time %v, yielding hash %v\n", nonce, count, diff, hash.Hex())
			ok = true
			return
		}

		nonce += 1
	}

	// Failed to find an appropriate nonce
	ok = false
	return
}
