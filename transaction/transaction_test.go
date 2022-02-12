package transaction

import (
	"testing"

	"github.com/AndrewCLu/TestcoinNode/common"
	"github.com/AndrewCLu/TestcoinNode/util"
)

// Tests that converting a transaction into a byte array and back yields the same transaction
func TestTransactionToByteArray(t *testing.T) {
	address := [common.AddressLength]byte{2, 3}
	amount := util.Float64UnitToUnit64Unit(69.69)

	output := TransactionOutput{ReceiverAddress: address, Amount: amount}
	transaction, success := New(
		[]*TransactionInput{},
		[]*TransactionOutput{&output},
	)

	if !success {
		t.Fatalf(`Failed to create a new coinbase transaction.`)
	}

	transactionBytes := transaction.Bytes()

	decodedTransaction := BytesToTransaction(transactionBytes)

	if !transaction.Equal(decodedTransaction) {
		t.Fatalf(`Decoded transaction is not equal to original. Original: %v, Decoded: %v`, transaction, decodedTransaction)
	}
}
