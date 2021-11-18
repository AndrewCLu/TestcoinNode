package account

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/protocol"
)

type account struct {
	address           [protocol.AddressLength]byte
	encodedPublicKey  []byte
	encodedPrivateKey []byte
}

func NewAccount() account {
	// TODO: Must resize length of address if AddressLength and HashLength are not the same
	address, encodedPublicKey, encodedPrivateKey := crypto.NewDigitalSignatureKeys()
	account := account{
		address:           address,
		encodedPublicKey:  encodedPublicKey,
		encodedPrivateKey: encodedPrivateKey,
	}

	fmt.Printf("Created account with address: %v\n", account.address)

	return account
}

func (a account) GetAddress() [protocol.AddressLength]byte {
	return a.address
}
