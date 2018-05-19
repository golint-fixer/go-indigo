package deprecated

import (
	"indogo/src/core/types"
	"math/big"
)

// Nonce - hash symbolizing amount of processing power required to produce block
type Nonce [8]byte

// Header -
type Header struct {
	ParentHash   types.Hash    `json:"parentHash"       gencodec:"required"`
	MinerAddress types.Address `json:"miner"            gencodec:"required"`
	TxHash       types.Hash    `json:"transactionsRoot" gencodec:"required"`
	Difficulty   *big.Int      `json:"difficulty"       gencodec:"required"`
	Time         *big.Int      `json:"timestamp"        gencodec:"required"`
	Extra        []byte        `json:"extraData"        gencodec:"required"`
	Nonce        Nonce         `json:"nonce"            gencodec:"required"`
}

// Body - container representing data of block
type Body struct {
	Transactions []*types.Transaction
	Uncles       []*Header
}

/*
// Block - container holding transactions in blockchain, may execute on filled condition
type Block struct {
	witness      *witness.Witness
	transactions Transactions

	hash atomic.Value
	size atomic.Value

	totaldifficulty *big.Int

	ReceivedAt   time.Time
	ReceivedFrom interface{}
}

// WithBody - returns a new block with given transactions.
func (b *Block) WithBody(transactions []*Transaction, uncles []*Header) *Block {
	block := &Block{
		header:       CopyHeader(b.header),
		transactions: make([]*Transaction, len(transactions)),
		uncles:       make([]*Header, len(uncles)),
	}
	copy(block.transactions, transactions)
	for i := range uncles {
		block.uncles[i] = CopyHeader(uncles[i])
	}
	return block
}
*/
