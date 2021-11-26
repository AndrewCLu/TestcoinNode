package node

import (
	"testing"

	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/transaction"
)

func TestSignVerifyTransactionInput(t *testing.T) {
	act := account.NewAccount()
	tx := NewCoinbaseTransaction(act, 10)

	outputPointer := transaction.TransactionOutputPointer{
		TransactionHash: tx.Hash(),
		OutputIndex:     uint16(0),
	}

	signature := SignInput(act.GetPrivateKey(), outputPointer)
	verified := VerifyInput(act.GetPublicKey(), outputPointer, signature)

	if !verified {
		t.Fatalf(`Failed to verify the signature of a new coinbase transaction.`)
	}
}
