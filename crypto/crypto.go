package crypto

import (
	"crypto/sha256"
)

const HashLength = 32

// Generates keys for an account, including a 32 byte address, public key, and private key
func NewAccountKeys() (address [HashLength]byte, encodedPublicKey []byte, encodedPrivateKey []byte) {
	publicKey, privateKey, _ := newECDSAKeyPair()

	encodedPublicKey, _ = encodePublicKey(publicKey)
	encodedPrivateKey, _ = encodePrivateKey(privateKey)
	address = HashBytes(encodedPublicKey)

	return address, encodedPublicKey, encodedPrivateKey
}

// Hashes bytes using SHA256
func HashBytes(bytes []byte) (hash [HashLength]byte) {
	return sha256.Sum256(bytes)
}
