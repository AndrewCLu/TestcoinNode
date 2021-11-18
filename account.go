package main

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/crypto"
)

const AddressLength = 32

type account struct {
	address           [AddressLength]byte
	encodedPublicKey  []byte
	encodedPrivateKey []byte
}

func NewAccount() account {
	// TODO: Must resize length of address if AddressLength and HashLength are not the same
	address, encodedPublicKey, encodedPrivateKey := crypto.NewAccountKeys()
	account := account{
		address:           address,
		encodedPublicKey:  encodedPublicKey,
		encodedPrivateKey: encodedPrivateKey,
	}

	fmt.Printf("Created account with address: %v\n", account.address)

	return account
}

func (a account) GetAddress() [AddressLength]byte {
	return a.address
}
