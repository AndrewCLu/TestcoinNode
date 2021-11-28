package crypto

import (
	"testing"
)

// Tests signing and verification using newly generated ECDSA keys
func TestSignVerify(t *testing.T) {
	publicKey, privateKey, _ := newECDSAKeyPair()
	encodedPublicKey, _ := encodePublicKey(publicKey)
	encodedPrivateKey, _ := encodePrivateKey(privateKey)

	bytesA := []byte("Go BTC")
	bytesB := []byte("Go ETH")

	s1, ok := SignByteArray(bytesA, encodedPrivateKey)
	if !ok {
		t.Fatalf(`Failed to sign array`)
	}

	verify1 := VerifyByteArray(bytesA, encodedPublicKey, s1)
	if !verify1 {
		t.Fatalf(`Verification false for byte array: %v`, bytesA)
	}

	s2, ok := SignByteArray(bytesA, encodedPrivateKey)
	if !ok {
		t.Fatalf(`Failed to sign array`)
	}
	verify2 := VerifyByteArray(bytesB, encodedPublicKey, s2)
	if verify2 {
		t.Fatalf(`Verification succeeded for two unequal byte arrays: %v and %v`, bytesA, bytesB)
	}
}
