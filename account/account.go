package account

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/protocol"
)

type Account struct {
	address           [protocol.AddressLength]byte
	encodedPublicKey  []byte
	encodedPrivateKey []byte
}

func NewAccount() Account {
	// TODO: Must resize length of address if AddressLength and HashLength are not the same
	address, encodedPublicKey, encodedPrivateKey := crypto.NewDigitalSignatureKeys()
	account := Account{
		address:           address,
		encodedPublicKey:  encodedPublicKey,
		encodedPrivateKey: encodedPrivateKey,
	}

	fmt.Printf("Created account with address: %v\n", account.address)

	return account
}

func (a Account) GetAddress() [protocol.AddressLength]byte {
	return a.address
}
