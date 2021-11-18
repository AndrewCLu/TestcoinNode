package util

import (
	"encoding/binary"
	"math"
)

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
func Uint16ToBytes(value uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, value)

	return bytes
}

// TODO: Safety checks on inputs
func BytesToUint16(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes)
}
