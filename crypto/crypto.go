package crypto

import (
	"crypto/sha256"
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/common"
)

// Generates digital signature, including a 32 byte address, public key, and private key
func NewDigitalSignatureKeys() (encodedPublicKey []byte, encodedPrivateKey []byte, err error) {
	publicKey, privateKey, err := newECDSAKeyPair()
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	// Encode keys using x509
	encodedPublicKey, err = encodePublicKey(publicKey)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	encodedPrivateKey, err = encodePrivateKey(privateKey)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	return encodedPublicKey, encodedPrivateKey, nil
}

// Hashes bytes using SHA256 and returns the corresponging Hash
func HashBytes(bytes []byte) common.Hash {
	hashBytes := sha256.Sum256(bytes)
	return common.BytesToHash(hashBytes[:])
}
