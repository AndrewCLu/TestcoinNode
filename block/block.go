package block

import (
	"bytes"
	"fmt"
	"math/rand"
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
	Target              [crypto.HashLength]byte `json:"target"`
	Nonce               uint32                  `json:"nonce"`
}

func (header BlockHeader) Solve() uint32 {
	var nonce uint32 = rand.Uint32()
	target := header.Target[:]

	for true {
		header.Nonce = nonce
		hash := header.Hash()
		fmt.Printf("Trying nonce %v, yielding hash %v", nonce, hash)

		if bytes.Compare(hash[:], target) < 0 {
			fmt.Printf("Successfully found nonce %v, yielding hash %v", nonce, hash)
			return nonce
		}

		nonce += 1
	}

	return nonce
}

func (header BlockHeader) BlockHeaderToByteArray() []byte {
	versionBytes := util.Uint16ToBytes(header.ProtocolVersion)

	previousBlockHashBytes := header.PreviousBlockHash[:]

	allTransactionsHashBytes := header.AllTransactionsHash[:]

	timeBytes, err := header.Timestamp.MarshalBinary()
	if err != nil {
		fmt.Printf("Error occurred creating byte array for transaction timestamp: %v\n", err)
	}

	targetBytes := header.Target[:]

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
