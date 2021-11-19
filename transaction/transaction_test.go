package transaction

import (
	"testing"

	"github.com/AndrewCLu/TestcoinNode/protocol"
	"github.com/AndrewCLu/TestcoinNode/util"
)

// Tests that converting a transaction into a byte array and back yields the same transaction
func TestTransactionToByteArray(t *testing.T) {
	address := [protocol.AddressLength]byte{2, 3}

	transaction, _ := NewCoinbaseTransaction(address, util.Float64UnitToUnit64Unit(69.69))

	transactionBytes := transaction.TransactionToByteArray()

	decodedTransaction := ByteArrayToTransaction(transactionBytes)

	if !AreTransactionsEqual(transaction, decodedTransaction) {
		t.Fatalf(`Decoded transaction is not equal to original. Original: %v, Decoded: %v`, transaction, decodedTransaction)
	}
}
