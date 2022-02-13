package protocol

import (
	"github.com/AndrewCLu/TestcoinNode/common"
)

const ProtocolVersionLength = 2  // Number of bytes used to denote the protocol version
const CurrentProtocolVersion = 1 // Current protocol version

const TestcoinUnitMultipler = 1000000000 // Actual account values are 1000000000 times less than the transaction amount values

const MaxTransactionsInBlock = 2 // How many transactions are allowed in each block

// Given the current block number, returns the appropriate target for solving proof of work
func ComputeTarget(blockNumber int) common.Target {
	// return [4]byte{0, 0, 0, 15}
	return [4]byte{0, 255, 255, 255}
}

// Given the current block number, return the appropriate coinbase reward for mining a block
func ComputeBlockReward(blockNumber int) uint64 {
	return 10
}
