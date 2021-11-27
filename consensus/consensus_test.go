package consensus

import (
	"testing"

	"github.com/AndrewCLu/TestcoinNode/account"
	"github.com/AndrewCLu/TestcoinNode/transaction"
	"github.com/AndrewCLu/TestcoinNode/util"
)

func TestSignVerifyTransactionInput(t *testing.T) {

	act := account.NewAccount()
	amount := util.Float64UnitToUnit64Unit(10.0)

	output := transaction.TransactionOutput{ReceiverAddress: act.GetAddress(), Amount: amount}
	tx, success := transaction.NewTransaction(
		[]transaction.TransactionInput{},
		[]transaction.TransactionOutput{output},
	)

	if !success {
		t.Fatalf(`Tried to make a new transaction but failed.`)
	}

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
