package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
)

func NewAccountKeys() (address [32]byte, encodedPublicKey []byte, encodedPrivateKey []byte) {
	publicKeyCurve := elliptic.P256()

	privateKey := new(ecdsa.PrivateKey)
	privateKey, err := ecdsa.GenerateKey(publicKeyCurve, rand.Reader)

	if err != nil {
		fmt.Println(err)
	}

	publicKey := &privateKey.PublicKey

	encodedPrivateKey, _ = x509.MarshalECPrivateKey(privateKey)
	encodedPublicKey, _ = x509.MarshalPKIXPublicKey(publicKey)
	address = sha256.Sum256(encodedPublicKey)

	fmt.Println("Private Key :")
	fmt.Printf("%v \n", encodedPrivateKey)
	fmt.Printf("%T \n", encodedPrivateKey)

	fmt.Println("Public Key :")
	fmt.Printf("%v \n", encodedPublicKey)
	fmt.Printf("%T \n", encodedPublicKey)

	fmt.Printf("%v \n", len(encodedPublicKey))
	fmt.Printf("%v \n", len(address))
	fmt.Printf("%v \n", address)

	return address, encodedPublicKey, encodedPrivateKey
}
