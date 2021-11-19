package util

import (
	"encoding/binary"
	"math"

	"github.com/AndrewCLu/TestcoinNode/protocol"
)

// TODO: Handle special cases where values are very large
func Float64UnitToUnit64Unit(value float64) uint64 {
	return uint64(value * protocol.TestcoinUnitMultipler)
}

// TODO: Handle special cases where values are very large
func Uint64UnitToFloat64Unit(value uint64) float64 {
	return float64(value) / protocol.TestcoinUnitMultipler
}

// TODO: Safety checks on inputs
func Float64ToBytes(value float64) []byte {
	bits := math.Float64bits(value)
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, bits)

	return bytes
}

// TODO: Safety checks on inputs
func BytesToFloat64(bytes []byte) float64 {
	bits := binary.BigEndian.Uint64(bytes)
	float := math.Float64frombits(bits)

	return float
}

// TODO: Safety checks on inputs
func Uint64ToBytes(value uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, value)

	return bytes
}

// TODO: Safety checks on inputs
func BytesToUint64(bytes []byte) uint64 {
	return binary.BigEndian.Uint64(bytes)
}

// TODO: Safety checks on inputs
func Uint16ToBytes(value uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, value)

	return bytes
}

// TODO: Safety checks on inputs
func BytesToUint16(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes)
}
