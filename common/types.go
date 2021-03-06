package common

import (
	"encoding/hex"
	"reflect"
)

const (
	HashLength    = 32 // The length of a SHA256 hash
	AddressLength = 32 // The length of a Testcoin address
	TargetLength  = 4  // A target represents the first 4 bytes of the hash a block must compare itself to
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

// Returns true if two hashes are equal else false
func (h1 Hash) Equal(h2 Hash) bool {
	return reflect.DeepEqual(h1, h2)
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

// Returns true if two hashes are equal else false
func (a1 Address) Equal(a2 Address) bool {
	return reflect.DeepEqual(a1, a2)
}

// A target is the first 4 bytes of the hash a block must compare itself to while mining.
type Target [TargetLength]byte

func BytesToTarget(bytes []byte) Target {
	var t Target

	if len(bytes) > TargetLength {
		bytes = bytes[:TargetLength]
	}

	copy(t[:], bytes)

	return t
}

// Converts a target into bytes.
func (t Target) Bytes() []byte {
	return t[:]
}

// Returns the full hash of a target for block hashes to compare themselves to
func (t Target) FullHash() Hash {
	var full [HashLength]byte

	for i := 0; i < TargetLength; i++ {
		full[i] = t[i]
	}

	// Fills rest of hash with zeros
	for i := TargetLength; i < HashLength; i++ {
		full[i] = byte(255)
	}

	return full
}
