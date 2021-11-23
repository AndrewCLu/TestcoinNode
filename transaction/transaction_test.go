package transaction

import (
	"testing"

	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/protocol"
	"github.com/AndrewCLu/TestcoinNode/util"
)

// Tests that converting a transaction into a byte array and back yields the same transaction
func TestTransactionToByteArray(t *testing.T) {
	address := [protocol.AddressLength]byte{2, 3}

	transaction, _ := NewCoinbaseTransaction(address, util.Float64UnitToUnit64Unit(69.69))

	transactionBytes := transaction.TransactionToByteArray()

	decodedTransaction := ByteArrayToTransaction(transactionBytes)

	if !transaction.Equal(decodedTransaction) {
		t.Fatalf(`Decoded transaction is not equal to original. Original: %v, Decoded: %v`, transaction, decodedTransaction)
	}
}

func TestSignVerifyTransactionInput(t *testing.T) {
	account := account.NewAccount()
	transaction, success := NewCoinbaseTransaction(account.GetAddress(), 10)

	if !success {
		t.Fatalf(`Failed to create a new coinbase transaction.`)
	}

	signature := SignInput(account.GetPrivateKey(), transaction.Hash(), uint16(0))
	verified := VerifyInput(account.GetPublicKey(), transaction.Hash(), uint16(0), signature)

	if !verified {
		t.Fatalf(`Failed to verify the signature of a new coinbase transaction.`)
	}
}
