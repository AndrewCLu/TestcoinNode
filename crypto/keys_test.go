package crypto

import (
	"reflect"
	"testing"
)

// Tests encoding and decoding of a public ECDSA key
func TestEncodeDecodePublicKey(t *testing.T) {
	publicKey, _, _ := newECDSAKeyPair()
	encodedBytes, _ := encodePublicKey(publicKey)
	decodedPublicKey, _ := decodePublicKey(encodedBytes)

	if !publicKey.Equal(decodedPublicKey) {
		t.Fatalf(`Public keys do not match. Original: %v, Decoded: %v`, publicKey, decodedPublicKey)
	}
}

// Tests encoding and decoding of a private ECDSA key
func TestEncodeDecodePrivateKey(t *testing.T) {
	_, privateKey, _ := newECDSAKeyPair()
	encodedBytes, _ := encodePrivateKey(privateKey)
	decodedPrivateKey, _ := decodePrivateKey(encodedBytes)

	if !privateKey.Equal(decodedPrivateKey) {
		t.Fatalf(`Private keys do not match. Original: %v, Decoded: %v`, privateKey, decodedPrivateKey)
	}
}

// Tests that the hash of a new account address corresponds to its public key
func TestAddressIsHashOfPublicKey(t *testing.T) {
	address, encodedPublicKey, _ := NewAccountKeys()
	hashedPublicKey := hashBytes(encodedPublicKey)

	if !reflect.DeepEqual(address, hashedPublicKey) {
		t.Fatalf(`Hashes do not match. Address: %v, Hashed public key: %v`, address, hashedPublicKey)
	}
}
