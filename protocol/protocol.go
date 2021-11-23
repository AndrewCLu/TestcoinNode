package protocol

const ProtocolVersionLength = 2
const CurrentProtocolVersion = 1

const AddressLength = 32

const TestcoinUnitMultipler = 1000000000 // Actual account values are 1000000000 times less than the transaction amount values

// Given the current block number, returns the appropriate target for solving proof of work
func DifficultyTargetFormula(blockNumber int) (target uint32) {
	return uint32(4)
}
