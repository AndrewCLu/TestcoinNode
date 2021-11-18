package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
)

// Hashes and then signs a byte array using an encoded ECDSA private key
// Returns signature as r and s in byte arrays
func signByteArray(bytes []byte, encodedPrivateKey []byte) (rBytes []byte, sBytes []byte, err error) {
	privateKey, decodeErr := decodePrivateKey(encodedPrivateKey)
	if decodeErr != nil {
		fmt.Println(decodeErr)
		return nil, nil, decodeErr
	}

	hashedBytes := HashBytes(bytes)

	r, s, signErr := ecdsa.Sign(rand.Reader, privateKey, hashedBytes[:])
	if signErr != nil {
		fmt.Println(signErr)
		return nil, nil, signErr
	}

	return r.Bytes(), s.Bytes(), nil
}

// Verifies a signature using an encoded ECDSA public key
// Returns true if verified else false
func verifyByteArray(bytes []byte, encodedPublicKey []byte, rBytes []byte, sBytes []byte) (verified bool, err error) {
	publicKey, decodeErr := decodePublicKey(encodedPublicKey)
	if decodeErr != nil {
		fmt.Println(decodeErr)
		return false, decodeErr
	}

	r := new(big.Int).SetBytes(rBytes)
	s := new(big.Int).SetBytes(sBytes)

	hashedBytes := HashBytes(bytes)
	verified = ecdsa.Verify(publicKey, hashedBytes[:], r, s)

	return verified, nil
}
