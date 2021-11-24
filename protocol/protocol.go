package protocol

const ProtocolVersionLength = 2  // Number of bytes used to denote the protocol version
const CurrentProtocolVersion = 1 // Current protocol version

const AddressLength = 32 // Length of a Testcoin address

const TestcoinUnitMultipler = 1000000000 // Actual account values are 1000000000 times less than the transaction amount values

const MaxTransactionsInBlock = 10 // How many transactions are allowed in each block

const TargetLength = 4 // Number of bytes we can adjust in the target

// Given the current block number, returns the appropriate target for solving proof of work
func ComputeTarget(blockNumber int) (target [TargetLength]byte) {
	return [4]byte{0, 0, 0, 15}
}
