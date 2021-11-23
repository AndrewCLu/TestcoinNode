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

func (header BlockHeader) Solve() int64 {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	var nonce uint32 = r.Uint32()
	target := header.Target[:]

	t1 := time.Now()
	// fmt.Printf("Solving with target %v ...\n", util.HashToHexString(header.Target))
	count := 0

	for true {
		count += 1
		header.Nonce = nonce
		hash := header.Hash()
		// fmt.Printf("Trying nonce %v, yielding hash %v\n", nonce, util.HashToHexString(hash))

		if bytes.Compare(hash[:], target) < 0 {
			t2 := time.Now()
			diff := t2.Sub(t1)

			// fmt.Printf("Successfully found nonce %v with %v tries in time %v, yielding hash %v\n", nonce, count, diff, util.HashToHexString(hash))
			return int64(diff / time.Second)
		}

		nonce += 1
	}

	t3 := time.Now()
	diff2 := t3.Sub(t1)
	return int64(diff2 / time.Second)
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
