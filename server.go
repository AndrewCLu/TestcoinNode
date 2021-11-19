package main

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/node"
)

func main() {
	fmt.Printf("Starting server at port 8080\n")

	node.InitializeNode()

	bob := account.NewAccount()
	alice := account.NewAccount()

	node.NewCoinbaseTransaction(bob, 69.69)
	node.NewCoinbaseTransaction(bob, 10)

	node.NewPeerTransaction(bob, alice.GetAddress(), 69)

	node.GetReadableAccountValue(bob)
	node.GetReadableAccountValue(alice)
}
