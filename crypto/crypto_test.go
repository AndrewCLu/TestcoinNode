package crypto

import (
	"testing"
)

// Tests signing and verification using newly generated account keys
func TestSignVerifyWithAccountKeys(t *testing.T) {
	_, encodedPublicKey, encodedPrivateKey := NewAccountKeys()

	bytesA := []byte("Go BTC")
	bytesB := []byte("Go ETH")

	s1, _ := signByteArray(bytesA, encodedPrivateKey)
	verify1, _ := verifyByteArray(bytesA, encodedPublicKey, s1)
	if !verify1 {
		t.Fatalf(`Verification false for byte array: %v`, bytesA)
	}

	s2, _ := signByteArray(bytesA, encodedPrivateKey)
	verify2, _ := verifyByteArray(bytesB, encodedPublicKey, s2)
	if verify2 {
		t.Fatalf(`Verification succeeded for two unequal byte arrays: %v and %v`, bytesA, bytesB)
	}
}
