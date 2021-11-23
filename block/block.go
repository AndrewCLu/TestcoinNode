package block

import (
	"fmt"
	"time"

	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/transaction"
	"github.com/AndrewCLu/TestcoinNode/util"
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

func (header BlockHeader) BlockHeaderToByteArray() []byte {
	versionBytes := util.Uint16ToBytes(header.ProtocolVersion)

	previousBlockHashBytes := header.PreviousBlockHash[:]

	allTransactionsHashBytes := header.AllTransactionsHash[:]

	timeBytes, err := header.Timestamp.MarshalBinary()
	if err != nil {
		fmt.Printf("Error occurred creating byte array for transaction timestamp: %v\n", err)
	}

	targetBytes := util.Uint32ToBytes(header.Target)

	nonceBytes := util.Uint32ToBytes(header.Nonce)

	allBytes := [][]byte{
		versionBytes,
		previousBlockHashBytes,
		allTransactionsHashBytes,
		timeBytes,
		targetBytes,
		nonceBytes,
	}

	return util.ConcatByteSlices(allBytes)
}

func (b Block) Hash() [crypto.HashLength]byte {
	return b.Header.Hash()
}

func (header BlockHeader) Hash() [crypto.HashLength]byte {
	return crypto.HashBytes(header.BlockHeaderToByteArray())
}
