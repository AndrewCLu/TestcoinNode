package main

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/node"
)

func main() {
	fmt.Printf("Starting server at port 8080\n")

	node.InitializeNode()

	bob := node.NewAccount()
	alice := node.NewAccount()

	node.GetReadableAccountValue(bob)
	node.GetReadableAccountValue(alice)

	node.NewCoinbaseTransaction(bob, 69.69)
	node.NewCoinbaseTransaction(bob, 10)

	node.NewPeerTransaction(bob, alice.GetAddress(), 69)

	node.GetReadableAccountValue(bob)
	node.GetReadableAccountValue(alice)

	node.NewPeerTransaction(alice, bob.GetAddress(), 6)

	node.GetReadableAccountValue(bob)
	node.GetReadableAccountValue(alice)
}
