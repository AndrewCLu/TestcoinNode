package account

import (
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

	return account
}

func (a Account) GetAddress() [protocol.AddressLength]byte {
	return a.address
}

func (a Account) GetPublicKey() []byte {
	return a.encodedPublicKey
}

func (a Account) GetPrivateKey() []byte {
	return a.encodedPrivateKey
}
