package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
)

// Generates keys for an account, including a 32 byte address, public key, and private key
func NewAccountKeys() (address [32]byte, encodedPublicKey []byte, encodedPrivateKey []byte) {
	publicKey, privateKey, _ := newECDSAKeyPair()

	encodedPublicKey, _ = encodePublicKey(publicKey)
	encodedPrivateKey, _ = encodePrivateKey(privateKey)
	address = hashKey(encodedPublicKey)

	return address, encodedPublicKey, encodedPrivateKey
}

// Generates an ECDSA public private key pair
func newECDSAKeyPair() (publicKey *ecdsa.PublicKey, privateKey *ecdsa.PrivateKey, err error) {
	publicKeyCurve := elliptic.P256()

	privateKey, err = ecdsa.GenerateKey(publicKeyCurve, rand.Reader)

	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	publicKey = &privateKey.PublicKey

	return publicKey, privateKey, nil
}

// Encodes a ECDSA public key using x509 encoding
func encodePublicKey(key *ecdsa.PublicKey) (bytes []byte, err error) {
	bytes, err = x509.MarshalPKIXPublicKey(key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return bytes, nil
}

// Decodes an x509 encoded ECDSA public key
func decodePublicKey(bytes []byte) (key *ecdsa.PublicKey, err error) {
	genericPublicKey, err := x509.ParsePKIXPublicKey(bytes)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	key = genericPublicKey.(*ecdsa.PublicKey)

	return key, nil
}

// Encodes a ECDSA private key using x509 encoding
func encodePrivateKey(key *ecdsa.PrivateKey) (bytes []byte, err error) {
	bytes, err = x509.MarshalECPrivateKey(key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return bytes, nil
}

// Decodes an x509 encoded ECDSA private key
func decodePrivateKey(bytes []byte) (key *ecdsa.PrivateKey, err error) {
	key, err = x509.ParseECPrivateKey(bytes)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return key, nil
}

// Hashes a key using SHA256
func hashKey(encodedKey []byte) (hash [32]byte) {
	return sha256.Sum256(encodedKey)
}
