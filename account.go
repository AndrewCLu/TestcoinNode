package main

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/crypto"
)

type account struct {
	address           [32]byte
	encodedPublicKey  []byte
	encodedPrivateKey []byte
}

func NewAccount() account {
	address, encodedPublicKey, encodedPrivateKey := crypto.NewAccountKeys()
	account := account{
		address:           address,
		encodedPublicKey:  encodedPublicKey,
		encodedPrivateKey: encodedPrivateKey,
	}

	fmt.Printf("Created account with address: %v\n", account.address)

	return account
}

func (a account) GetAddress() [32]byte {
	return a.address
}
