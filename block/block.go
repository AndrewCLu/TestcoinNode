package block

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/protocol"
	"github.com/AndrewCLu/TestcoinNode/transaction"
	"github.com/AndrewCLu/TestcoinNode/util"
)

type Block struct {
	Header BlockHeader               `json:"header"`
	Body   []transaction.Transaction `json:"body"`
}

type BlockHeader struct {
	ProtocolVersion     uint16                      `json:"protocolVersion"`
	PreviousBlockHash   [crypto.HashLength]byte     `json:"previousBlockHash"`
	AllTransactionsHash [crypto.HashLength]byte     `json:"allTransactionsHash"`
	Timestamp           time.Time                   `json:"timestamp"`
	Target              [protocol.TargetLength]byte `json:"target"`
	Nonce               uint32                      `json:"nonce"`
}

func NewBlock(currentBlockNumber uint64, previousBlockHash [crypto.HashLength]byte, transactions []transaction.Transaction) (Block, error) {
	if len(transactions) > protocol.MaxTransactionsInBlock {
		return Block{}, errors.New("Number of transactions exceeds max allowable.")
	}

	var transactionHashes [][]byte
	for i, tx := range transactions {
		hash := tx.Hash()
		transactionHashes[i] = hash[:]
	}
	allTransactionsHash := crypto.HashBytes(util.ConcatByteSlices(transactionHashes))

	header := BlockHeader{
		ProtocolVersion:     protocol.CurrentProtocolVersion,
		PreviousBlockHash:   previousBlockHash,
		AllTransactionsHash: allTransactionsHash,
		Timestamp:           time.Now().Round(0),
		Target:              protocol.ComputeTarget(currentBlockNumber),
		Nonce:               uint32(0),
	}

	block := Block{
		Header: header,
		Body:   transactions,
	}

	return block, nil
}

// Given a block header, compute the nonce that results in a hash under the desired target
func (header BlockHeader) Solve() uint32 {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	var nonce uint32 = r.Uint32()

	var target [crypto.HashLength]byte
	for i := 0; i < protocol.TargetLength; i++ {
		target[i] = header.Target[i]
	}
	for i := protocol.TargetLength; i < crypto.HashLength; i++ {
		target[i] = byte(255)
	}
	targetBytes := target[:]

	t1 := time.Now()
	fmt.Printf("Solving with target %v ...\n", util.HashToHexString(target))

	count := 0
	for true {
		count += 1
		header.Nonce = nonce
		hash := header.Hash()

		// fmt.Printf("Trying nonce %v, yielding hash %v\n", nonce, util.HashToHexString(hash))

		if bytes.Compare(hash[:], targetBytes) < 0 {
			t2 := time.Now()
			diff := t2.Sub(t1)

			fmt.Printf("Successfully found nonce %v with %v tries in time %v, yielding hash %v\n", nonce, count, diff, util.HashToHexString(hash))
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
