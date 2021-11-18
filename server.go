package main

import (
	"fmt"
	// "log"
	// "net/http"
)

const CurrentProtocolVersion = 1

func main() {
	fmt.Printf("Starting server at port 8080\n")

	account := NewAccount()

	t := NewCoinbaseTransaction(account.GetAddress(), 69.69)
	fmt.Println(t.TransactionToByteArray())
}
