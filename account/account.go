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
	encodedPublicKey, encodedPrivateKey := crypto.NewDigitalSignatureKeys()
	address := GetAddressFromPublicKey(encodedPublicKey)

	account := Account{
		address:           address,
		encodedPublicKey:  encodedPublicKey,
		encodedPrivateKey: encodedPrivateKey,
	}

	return account
}

func GetAddressFromPublicKey(publicKey []byte) (address [protocol.AddressLength]byte) {
	return crypto.HashBytes(publicKey)
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
