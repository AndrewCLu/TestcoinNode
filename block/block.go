package block

import (
	"fmt"
	"time"

	"github.com/AndrewCLu/TestcoinNode/common"
	"github.com/AndrewCLu/TestcoinNode/crypto"
	"github.com/AndrewCLu/TestcoinNode/protocol"
	"github.com/AndrewCLu/TestcoinNode/transaction"
	"github.com/AndrewCLu/TestcoinNode/util"
)

// A block is a collection of transactions, including a coinbase, with a header containing metadata
type Block struct {
	Header   *BlockHeader               `json:"header"`
	Body     []*transaction.Transaction `json:"body"`
	Coinbase *transaction.Transaction   `json:"coinbase"`
}

// A block header describes metatdata of a block, including a hash of its transactions,
// the previous block hash, a mining difficulty target, and a nonce
type BlockHeader struct {
	ProtocolVersion     uint16        `json:"protocolVersion"`
	PreviousBlockHash   common.Hash   `json:"previousBlockHash"`
	AllTransactionsHash common.Hash   `json:"allTransactionsHash"`
	Timestamp           time.Time     `json:"timestamp"`
	Target              common.Target `json:"target"`
	Nonce               uint32        `json:"nonce"`
}

// Generates a new block given the state of the previous blocks and a list of transactions
// Returns a pointer to the block if valid block can be generated
func New(
	previousBlockHash common.Hash,
	currentBlockNumber int,
	transactions []*transaction.Transaction,
	coinbase *transaction.Transaction,
) (blk *Block, ok bool) {
	// Check number of transactions is under the limit
	if len(transactions) > protocol.MaxTransactionsInBlock {
		fmt.Println("Number of transactions exceeds max allowable.")
		return nil, false
	}

	// Hash all the transactions, including coinbase
	transactionHashes := make([][]byte, len(transactions)+1)
	for i, tx := range transactions {
		hash := tx.Hash()
		transactionHashes[i] = hash.Bytes()
	}
	transactionHashes[len(transactions)] = coinbase.Hash().Bytes()
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
		Header:   &header,
		Body:     transactions,
		Coinbase: coinbase,
	}

	return &block, true
}

// Converts a BlockHeader into byte representation
func (header *BlockHeader) Bytes() []byte {
	versionBytes := util.Uint16ToBytes(header.ProtocolVersion)

	previousBlockHashBytes := header.PreviousBlockHash.Bytes()

	allTransactionsHashBytes := header.AllTransactionsHash.Bytes()

	timeBytes, err := header.Timestamp.MarshalBinary()
	if err != nil {
		fmt.Printf("Error occurred creating byte array for transaction timestamp: %v\n", err)
	}

	targetBytes := header.Target.Bytes()

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

// Returns the hash of a block, which is simply the hash of the block header
func (b *Block) Hash() common.Hash {
	return b.Header.Hash()
}

// Returns the hash of a block header by hashing the byte representation of the header
func (header *BlockHeader) Hash() common.Hash {
	return crypto.HashBytes(header.Bytes())
}
