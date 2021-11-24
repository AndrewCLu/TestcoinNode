package block

import (
	"errors"
	"fmt"
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

func NewBlock(previousBlockHash [crypto.HashLength]byte, currentBlockNumber int, transactions []transaction.Transaction) (Block, error) {
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
