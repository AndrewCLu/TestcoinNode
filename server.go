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
	fmt.Printf("Address: %v\n", account.GetAddress())

	NewCoinbaseTransaction(account.GetAddress(), 69.69)
}
