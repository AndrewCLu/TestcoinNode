package miner

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/AndrewCLu/TestcoinNode/block"
	"github.com/AndrewCLu/TestcoinNode/chain"
	"github.com/AndrewCLu/TestcoinNode/protocol"
	"github.com/AndrewCLu/TestcoinNode/util"
)

// Tries to mine a block from pending transactions on the chain
// Returns the block and a boolean indicating success
// TODO: Have better selection criteria for transactions
func MineBlock() (blk *block.Block, ok bool) {
	txs, txOk := chain.GetPendingTransactions(protocol.MaxTransactionsInBlock)
	if !txOk {
		fmt.Printf("Could not get transactions from the current chain")
		return nil, false
	}

	lastBlockHash, lastBlockNum, lastBlockOk := chain.GetLastBlockInfo()
	if !lastBlockOk {
		fmt.Printf("Could not get last block from current chain")
		return nil, false
	}
	blockNum := lastBlockNum + 1

	block, blockOk := block.NewBlock(lastBlockHash, blockNum, txs)
	if !blockOk {
		fmt.Printf("Failed to create new block")
		return nil, false
	}

	// Compute the nonce that solves a block header
	nonce := Solve(block.Header)
	block.Header.Nonce = nonce

	blockHash := block.Hash()
	for _, tx := range block.Body {
		fmt.Printf("Block %v confirmed transaction %v\n", util.HashToHexString(blockHash), util.HashToHexString(tx.Hash()))
	}

	return *block
}

// Given a block header, compute the nonce that results in a hash under the desired target
// TODO: Copy the block header over before computing nonces
func Solve(header block.BlockHeader) uint32 {
	// Start the nonce at a random number to avoid multiple nodes mining the same nonces
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	var nonce uint32 = r.Uint32()

	target := protocol.GetFullTargetFromHeader(header.Target)

	t1 := time.Now()
	fmt.Printf("Solving block with target %v ...\n", util.HashToHexString(target))

	count := 0
	for true {
		count += 1
		header.Nonce = nonce
		hash := header.Hash()

		// fmt.Printf("Trying nonce %v, yielding hash %v\n", nonce, util.HashToHexString(hash))

		if bytes.Compare(hash[:], target[:]) < 0 {
			t2 := time.Now()
			diff := t2.Sub(t1)

			fmt.Printf("Successfully found nonce %v with %v tries in time %v, yielding hash %v\n", nonce, count, diff, util.HashToHexString(hash))
			return nonce
		}

		nonce += 1
	}

	return nonce
}
