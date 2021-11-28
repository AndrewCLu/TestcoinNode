package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
)

// An ECDSASignature is an elliptic curve cryptography signature which consists of two big integers r and s
type ECDSASignature struct {
	r big.Int
	s big.Int
}

// Converts a byte array into a ECDSASignature
// TODO: Find a better way to convert this that places the conversion info inside the struct
func BytesToECDSASignature(bytes []byte) ECDSASignature {
	rBytes := bytes[:len(bytes)/2]
	sBytes := bytes[len(bytes)/2:]
	rptr := new(big.Int).SetBytes(rBytes)
	sptr := new(big.Int).SetBytes(sBytes)

	return ECDSASignature{
		r: *rptr,
		s: *sptr,
	}
}

// Converts an ECDSASignature into bytes
func (sig ECDSASignature) Bytes() []byte {
	rBytes := sig.r.Bytes()
	sBytes := sig.s.Bytes()

	signatureBytes := make([]byte, len(rBytes), len(rBytes)+len(sBytes)) // Allocate extra capacity for appending sBytes
	copy(signatureBytes, rBytes)
	signatureBytes = append(signatureBytes, sBytes...)

	return signatureBytes
}

// Hashes and then signs a byte array using an encoded ECDSA private key
// Returns signature as (r concat s) in a byte array
func SignByteArray(bytes []byte, privateKey []byte) (signature *ECDSASignature, ok bool) {
	// Must decode private key first
	decodedPrivateKey, decodeErr := decodePrivateKey(privateKey)
	if decodeErr != nil {
		fmt.Println(decodeErr)
		return nil, false
	}

	hash := HashBytes(bytes)

	r, s, signErr := ecdsa.Sign(rand.Reader, decodedPrivateKey, hash.Bytes())
	if signErr != nil {
		fmt.Println(signErr)
		return nil, false
	}

	signature = &ECDSASignature{
		r: *r,
		s: *s,
	}

	return signature, true
}

// Verifies a signature using an encoded ECDSA public key
// Returns true if verified else false
func VerifyByteArray(bytes []byte, publicKey []byte, signature *ECDSASignature) (verified bool) {
	// Must decode public key first
	decodedPublicKey, decodeErr := decodePublicKey(publicKey)
	if decodeErr != nil {
		fmt.Println(decodeErr)
		return false
	}

	r := signature.r
	s := signature.s

	hash := HashBytes(bytes)
	verified = ecdsa.Verify(decodedPublicKey, hash.Bytes(), &r, &s)

	return verified
}
