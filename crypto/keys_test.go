package crypto

import (
	"reflect"
	"testing"
)

// Tests signing and verification using newly generated account keys
func TestSignVerifyWithAccountKeys(t *testing.T) {
	_, encodedPublicKey, encodedPrivateKey := NewAccountKeys()

	bytesA := []byte("Go BTC")
	bytesB := []byte("Go ETH")

	r1, s1, _ := signByteArray(bytesA, encodedPrivateKey)
	verify1, _ := verifyByteArray(bytesA, encodedPublicKey, r1, s1)
	if !verify1 {
		t.Fatalf(`Verification false for byte array: %v`, bytesA)
	}

	r2, s2, _ := signByteArray(bytesA, encodedPrivateKey)
	verify2, _ := verifyByteArray(bytesB, encodedPublicKey, r2, s2)
	if verify2 {
		t.Fatalf(`Verification succeeded for two unequal byte arrays: %v and %v`, bytesA, bytesB)
	}
}

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
