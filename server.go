package main

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/node"
)

func main() {
	fmt.Printf("Starting server at port 8080\n")

	node.InitializeNode()

	account := account.NewAccount()

	node.NewCoinbaseTransaction(account, 69.69)
	node.NewCoinbaseTransaction(account, 10)

	node.GetAccountValue(account)
}
