package crypto

import (
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
