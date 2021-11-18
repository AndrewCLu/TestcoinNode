package main

import (
	"fmt"

	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/transaction"
)

func main() {
	fmt.Printf("Starting server at port 8080\n")

	account := account.NewAccount()

	t := transaction.NewCoinbaseTransaction(account.GetAddress(), 69.69)
	fmt.Println(t.TransactionToByteArray())
}
