package transaction

import (
	"testing"

	"github.com/AndrewCLu/TestcoinNode/protocol"
)

// Tests that converting a transaction into a byte array and back yields the same transaction
func TestTransactionToByteArray(t *testing.T) {
	address := [protocol.AddressLength]byte{2, 3}

	transaction := NewCoinbaseTransaction(address, 69.69)

	transactionBytes := transaction.TransactionToByteArray()

	decodedTransaction := ByteArrayToTransaction(transactionBytes)

	if !AreTransactionsEqual(transaction, decodedTransaction) {
		t.Fatalf(`Decoded transaction is not equal to original. Original: %v, Decoded: %v`, transaction, decodedTransaction)
	}
}
