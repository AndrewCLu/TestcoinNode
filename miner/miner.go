package miner

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/AndrewCLu/TestcoinNode/block"
	"github.com/AndrewCLu/TestcoinNode/chain"
	"github.com/AndrewCLu/TestcoinNode/protocol"
)

// Tries to mine a block from pending transactions on the chain
// Returns the block and a boolean indicating success
// TODO: Have better selection criteria for transactions
func MineBlock() (blk *block.Block, ok bool) {
	txs, txOk := chain.GetPendingTransactions(protocol.MaxTransactionsInBlock)
	if !txOk {
		fmt.Println("Could not get transactions from the current chain")
		return nil, false
	}

	lastBlockHash, lastBlockNum, lastBlockOk := chain.GetLastBlockInfo()
	if !lastBlockOk {
		fmt.Println("Could not get last block from current chain")
		return nil, false
	}
	blockNum := lastBlockNum + 1

	block, blockOk := block.NewBlock(lastBlockHash, blockNum, txs)
	if !blockOk {
		fmt.Println("Failed to create new block")
		return nil, false
	}

	// Compute the nonce that solves a block header
	nonce := Solve(block.Header)
	block.Header.Nonce = nonce

	blockHash := block.Hash()
	for _, tx := range block.Body {
		fmt.Printf("Block %v confirmed transaction %v\n", blockHash.Hex(), tx.Hash().Hex())
	}

	return *block
}

// Given a block header, compute the nonce that results in a hash under the desired target
// Does not modify the block header passed in
// TOOD: Make sure Solve does not modify the original header
func Solve(header block.BlockHeader) uint32 {
	// Start the nonce at a random number to avoid multiple nodes mining the same nonces
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	var nonce uint32 = r.Uint32()

	// Represent target as a full hash to compare block hash to
	target := protocol.GetFullTargetFromHeader(header.Target)

	t1 := time.Now()
	fmt.Printf("Solving block with target %v ...\n", target.Hex())

	count := 0
	for true {
		count += 1
		header.Nonce = nonce
		hash := header.Hash()

		if bytes.Compare(hash.Bytes(), target.Bytes()) < 0 {
			t2 := time.Now()
			diff := t2.Sub(t1)

			fmt.Printf("Successfully found nonce %v with %v tries in time %v, yielding hash %v\n", nonce, count, diff, hash.Hex())
			return nonce
		}

		nonce += 1
	}

	return nonce
}
