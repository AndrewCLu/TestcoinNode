package protocol

const ProtocolVersionLength = 2
const CurrentProtocolVersion = 1

const AddressLength = 32

const TestcoinUnitMultipler = 1000000000 // Actual account values are 1000000000 times less than the transaction amount values

const TargetLength = 4

// Given the current block number, returns the appropriate target for solving proof of work
func TargetFormula(blockNumber int) (target [TargetLength]byte) {
	return [4]byte{0, 0, 0, 15}
}
