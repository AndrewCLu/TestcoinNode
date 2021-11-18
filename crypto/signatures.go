package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
)

const SignatureLength = 64

// Hashes and then signs a byte array using an encoded ECDSA private key
// Returns signature as r and s in byte arrays
func signByteArray(bytes []byte, encodedPrivateKey []byte) (signature [SignatureLength]byte, err error) {
	privateKey, decodeErr := decodePrivateKey(encodedPrivateKey)
	if decodeErr != nil {
		fmt.Println(decodeErr)
		return [SignatureLength]byte{}, decodeErr
	}

	hashedBytes := HashBytes(bytes)

	r, s, signErr := ecdsa.Sign(rand.Reader, privateKey, hashedBytes[:])
	if signErr != nil {
		fmt.Println(signErr)
		return [SignatureLength]byte{}, signErr
	}

	signatureBytes := append(r.Bytes(), s.Bytes()...)
	copy(signature[:], signatureBytes)

	// fmt.Printf("signing %v with key %v, signature is: %v\n", bytes, encodedPrivateKey, signature)
	return signature, nil
}

// Verifies a signature using an encoded ECDSA public key
// Returns true if verified else false
func verifyByteArray(bytes []byte, encodedPublicKey []byte, signature [SignatureLength]byte) (verified bool, err error) {
	publicKey, decodeErr := decodePublicKey(encodedPublicKey)
	if decodeErr != nil {
		fmt.Println(decodeErr)
		return false, decodeErr
	}

	rBytes := signature[:(SignatureLength / 2)]
	sBytes := signature[(SignatureLength / 2):]
	r := new(big.Int).SetBytes(rBytes)
	s := new(big.Int).SetBytes(sBytes)

	hashedBytes := HashBytes(bytes)
	verified = ecdsa.Verify(publicKey, hashedBytes[:], r, s)

	// fmt.Printf("verifying %v with signature %v, verification is: %v\n", bytes, signature, verified)
	return verified, nil
}
