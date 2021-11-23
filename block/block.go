package block

import (
	"time"

	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/transaction"
)

type Block struct {
	Header BlockHeader               `json:"header"`
	Body   []transaction.Transaction `json:"body"`
}

type BlockHeader struct {
	ProtocolVersion     uint16                  `json:"protocolVersion"`
	PreviousBlockHash   [crypto.HashLength]byte `json:"previousBlockHash"`
	AllTransactionsHash [crypto.HashLength]byte `json:"allTransactionsHash"`
	Timestamp           time.Time               `json:"timestamp"`
	Target              uint32                  `json:"target"`
	Nonce               uint32                  `json:"nonce"`
}
