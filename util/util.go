package util

import (
	"encoding/binary"
	"math"

	"github.com/AndrewCLu/TestcoinNode/protocol"
)

// Concatenates a slice of byte slices into one slice
func ConcatByteSlices(slices [][]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}

	allBytes := make([]byte, totalLen)

	var i int
	for _, s := range slices {
		i += copy(allBytes[i:], s)
	}

	return allBytes
}

// TODO: Handle special cases where values are very large
// Converts a human readable float representation of a Testcoin amount into a transactional amount
func Float64UnitToUnit64Unit(value float64) uint64 {
	return uint64(value * protocol.TestcoinUnitMultiplier)
}

// TODO: Handle special cases where values are very large
// Converts a uin64 transactional representation of a Testcoin amount innto human readable float units
func Uint64UnitToFloat64Unit(value uint64) float64 {
	return float64(value) / protocol.TestcoinUnitMultiplier
}

// Converts a uint16 to byte representation
func Uint16ToBytes(value uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, value)

	return bytes
}

// Converts uint16 as bytes to a uint16
func BytesToUint16(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes)
}

// Converts a uint32 to byte representation
func Uint32ToBytes(value uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, value)

	return bytes
}

// Converts uint32 as bytes to a uint32
func BytesToUint32(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
}

// Converts a uint64 to byte representation
func Uint64ToBytes(value uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, value)

	return bytes
}

// Converts uint64 as bytes to a uint64
func BytesToUint64(bytes []byte) uint64 {
	return binary.BigEndian.Uint64(bytes)
}

// Converts a float64 to byte representation
func Float64ToBytes(value float64) []byte {
	bits := math.Float64bits(value)
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, bits)

	return bytes
}

// Converts float64 as bytes to a float64
func BytesToFloat64(bytes []byte) float64 {
	bits := binary.BigEndian.Uint64(bytes)
	float := math.Float64frombits(bits)

	return float
}
