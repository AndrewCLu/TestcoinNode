package pow

import (
	"testing"

	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/transaction"
	"github.com/AndrewCLu/TestcoinNode/util"
)

func TestSignVerifyTransactionInput(t *testing.T) {
	pow, _ := New()

	act, _ := account.New()
	amount := util.Float64UnitToUnit64Unit(10.0)

	output := &transaction.TransactionOutput{ReceiverAddress: act.Address, Amount: amount}
	tx, success := transaction.New(
		[]*transaction.TransactionInput{},
		[]*transaction.TransactionOutput{output},
	)

	if !success {
		t.Fatalf(`Tried to make a new transaction but failed.`)
	}

	outputPointer := &transaction.TransactionOutputPointer{
		TransactionHash: tx.Hash(),
		OutputIndex:     uint16(0),
	}

	signature := pow.SignInput(act.PrivateKey, outputPointer)
	verified := pow.VerifyInput(act.PublicKey, outputPointer, signature)

	if !verified {
		t.Fatalf(`Failed to verify the signature of a new coinbase transaction.`)
	}
}
