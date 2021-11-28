package account

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/common"
	"github.com/AndrewCLu/TestcoinNode/crypto"
)

// An account represents a Testcoin user that is using this node as their wallet
type Account struct {
	Address    common.Address // The public address of the wallet
	PublicKey  []byte         // The x509 encoded public key of the wallet
	PrivateKey []byte         // The x509 encoded private key of the wallet
}

// Generates a new account with a new pair of public and private keys
// Returns a boolean indicating success
func NewAccount() (*Account, bool) {
	encodedPublicKey, encodedPrivateKey, err := crypto.NewDigitalSignatureKeys()
	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	address := GetAddressFromPublicKey(encodedPublicKey)

	account := Account{
		Address:    address,
		PublicKey:  encodedPublicKey,
		PrivateKey: encodedPrivateKey,
	}

	return &account, true
}

// Given a public key, return the corresponding address
func GetAddressFromPublicKey(publicKey []byte) common.Address {
	return common.BytesToAddress(crypto.HashBytes(publicKey).Bytes())
}
