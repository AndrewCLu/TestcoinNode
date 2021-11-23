package util

import (
	"encoding/binary"
	"encoding/hex"
	"math"

	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/protocol"
)

func AddressToHexString(bytes [protocol.AddressLength]byte) string {
	return hex.EncodeToString(bytes[:])
}

func HexStringToAddress(str string) [protocol.AddressLength]byte {
	var address [protocol.AddressLength]byte
	decodedBytes, _ := hex.DecodeString(str)

	copy(address[:], decodedBytes)

	return address
}

func HashToHexString(bytes [crypto.HashLength]byte) string {
	return hex.EncodeToString(bytes[:])
}

func HexStringToHash(str string) [crypto.HashLength]byte {
	var hash [crypto.HashLength]byte
	decodedBytes, _ := hex.DecodeString(str)

	copy(hash[:], decodedBytes)

	return hash
}

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
func Float64UnitToUnit64Unit(value float64) uint64 {
	return uint64(value * protocol.TestcoinUnitMultipler)
}

// TODO: Handle special cases where values are very large
func Uint64UnitToFloat64Unit(value uint64) float64 {
	return float64(value) / protocol.TestcoinUnitMultipler
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

// TODO: Safety checks on inputs
func Uint32ToBytes(value uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, value)

	return bytes
}

// TODO: Safety checks on inputs
func BytesToUint32(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
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
