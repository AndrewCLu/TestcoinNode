package protocol

import "github.com/AndrewCLu/TestcoinNode/crypto"

const ProtocolVersionLength = 2  // Number of bytes used to denote the protocol version
const CurrentProtocolVersion = 1 // Current protocol version

// const AddressLength = 32 // Length of a Testcoin address

const TestcoinUnitMultipler = 1000000000 // Actual account values are 1000000000 times less than the transaction amount values

const MaxTransactionsInBlock = 10 // How many transactions are allowed in each block

const TargetLength = 4 // Number of bytes we can adjust in the target

// Given the current block number, returns the appropriate target for solving proof of work
func ComputeTarget(blockNumber int) (target [TargetLength]byte) {
	// return [4]byte{0, 0, 0, 15}
	return [4]byte{0, 255, 255, 255}
}

// Given the first few bytes of the target stored in the block header, return the full target as a hash to be compared to
func GetFullTargetFromHeader(targetHeader [TargetLength]byte) [crypto.HashLength]byte {
	var target [crypto.HashLength]byte
	for i := 0; i < TargetLength; i++ {
		target[i] = targetHeader[i]
	}
	for i := TargetLength; i < crypto.HashLength; i++ {
		target[i] = byte(255)
	}

	return target
}
