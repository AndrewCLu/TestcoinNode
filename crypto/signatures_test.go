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
