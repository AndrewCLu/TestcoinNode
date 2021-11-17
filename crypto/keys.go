package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"os"
)

func GenKey() {
	pubkeyCurve := elliptic.P256() //see http://golang.org/pkg/crypto/elliptic/#P256

	privatekey := new(ecdsa.PrivateKey)
	privatekey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader) // this generates a public & private key pair

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var pubkey ecdsa.PublicKey
	pubkey = privatekey.PublicKey

	x509EncodedPriv, _ := x509.MarshalECPrivateKey(privatekey)
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(&pubkey)

	fmt.Println("Private Key :")
	fmt.Printf("%x \n", x509EncodedPriv)
	fmt.Printf("%T \n", x509EncodedPriv)

	fmt.Println("Public Key :")
	fmt.Printf("%x \n", x509EncodedPub)
	fmt.Printf("%T \n", x509EncodedPub)

	fmt.Printf("%v \n", len(x509EncodedPub))
	fmt.Printf("%v \n", len(sha256.Sum256(x509EncodedPub)))
}
