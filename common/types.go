package common

import (
	"encoding/hex"
)

const (
	HashLength    = 32
	AddressLength = 32
)

// A hash is a 32 byte SHA256 hash of data.
type Hash [HashLength]byte

// Given some bytes, return a Hash represented by the bytes.
// Trims bytes if length exceeds that of the Hash type.
func BytesToHash(bytes []byte) Hash {
	var h Hash

	if len(bytes) > HashLength {
		bytes = bytes[:HashLength]
	}

	copy(h[:], bytes)

	return h
}

// Gets the byte representation of the hash.
func (h Hash) Bytes() []byte {
	return h[:]
}

// Converts a hash into a hex string.
func (h Hash) Hex() string {
	return hex.EncodeToString(h.Bytes())
}

// An address is a 32 byte address of a Testcoin account.
type Address [AddressLength]byte

// Given some bytes, return an Address represented by the bytes.
// Trims bytes if length exceeds that of the Address type.
func BytesToAddress(bytes []byte) Address {
	var a Address

	if len(bytes) > AddressLength {
		bytes = bytes[:AddressLength]
	}

	copy(a[:], bytes)

	return a
}

// Gets the byte representation of the address.
func (a Address) Bytes() []byte {
	return a[:]
}

// Converts an address into a hex string.
func (a Address) Hex() string {
	return hex.EncodeToString(a.Bytes())
}
