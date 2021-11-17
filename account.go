package main

import "fmt"

type account struct {
	publicAddress string
	privateKey    string
}

func NewAccount() account {
	account := account{publicAddress: "asdasd", privateKey: "b"}
	fmt.Printf("Created account with address: %v\n", account.publicAddress)

	return account
}

func (a account) GetAddress() string {
	return a.publicAddress
}
