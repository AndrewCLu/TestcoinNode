package crypto

import (
	"crypto/sha256"
)

const HashLength = 32

// Generates digital signature, including a 32 byte address, public key, and private key
func NewDigitalSignatureKeys() (address [HashLength]byte, encodedPublicKey []byte, encodedPrivateKey []byte) {
	publicKey, privateKey, _ := newECDSAKeyPair()

	encodedPublicKey, _ = encodePublicKey(publicKey)
	encodedPrivateKey, _ = encodePrivateKey(privateKey)
	address = GetAddressFromPublicKey(encodedPublicKey)

	return address, encodedPublicKey, encodedPrivateKey
}

func GetAddressFromPublicKey(publicKey []byte) (address [HashLength]byte) {
	return HashBytes(publicKey)
}

// Hashes bytes using SHA256
func HashBytes(bytes []byte) (hash [HashLength]byte) {
	return sha256.Sum256(bytes)
}
