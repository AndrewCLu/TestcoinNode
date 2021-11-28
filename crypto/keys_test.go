package crypto

import (
	"testing"
)

// Tests encoding and decoding of a public ECDSA key
func TestEncodeDecodePublicKey(t *testing.T) {
	publicKey, _, err := newECDSAKeyPair()
	if err != nil {
		t.Fatalf(`Failed to create new digital signature keys`)
	}

	encodedBytes, err := encodePublicKey(publicKey)
	if err != nil {
		t.Fatalf(`Failed to encode public key`)
	}

	decodedPublicKey, err := decodePublicKey(encodedBytes)
	if err != nil {
		t.Fatalf(`Failed to decode public key`)
	}

	if !publicKey.Equal(decodedPublicKey) {
		t.Fatalf(`Public keys do not match. Original: %v, Decoded: %v`, publicKey, decodedPublicKey)
	}
}

// Tests encoding and decoding of a private ECDSA key
func TestEncodeDecodePrivateKey(t *testing.T) {
	_, privateKey, err := newECDSAKeyPair()
	if err != nil {
		t.Fatalf(`Failed to create new digital signature keys`)
	}

	encodedBytes, err := encodePrivateKey(privateKey)
	if err != nil {
		t.Fatalf(`Failed to encode private key`)
	}

	decodedPrivateKey, err := decodePrivateKey(encodedBytes)
	if err != nil {
		t.Fatalf(`Failed to decode private key`)
	}

	if !privateKey.Equal(decodedPrivateKey) {
		t.Fatalf(`Private keys do not match. Original: %v, Decoded: %v`, privateKey, decodedPrivateKey)
	}
}
