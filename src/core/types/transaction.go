package types

import (
	contracts "indogo/indo-go/src/contracts"
	"indogo/indo-go/src/witness"
	"math/big"
	"sync/atomic"
	"time"
)

//Transaction - Data representing transfer of value (can be null), as well as the transfer of data via payload. May be triggered on conditions, set via smart contract.
type Transaction struct {
	Data transactiondata `json:"txdata"`

	Contract *contracts.Contract `json:"contract"`

	verifications *big.Int
	weight        *big.Int

	initialwitness *witness.Witness

	ReceivedAt   time.Time
	ReceivedFrom interface{}

	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

type transactiondata struct {
	// Initialized in func:
	Nonce     uint64    `json:"nonce" gencodec:"required"`
	Recipient *Address  `json:"recipient" gencodec:"required"`
	Amount    *big.Int  `json:"value" gencodec:"required"`
	Payload   []byte    `json:"payload" gencodec:"required"`
	Time      time.Time `json:"timestamp" gencodec:"required"`
	Extra     []byte    `json:"extraData" gencodec:"required"`

	// Initialized at intercept:
	Hash       *Hash `json:"hash" gencodec:"required"`
	ParentHash *Hash `json:"parentHash" gencodec:"required"`
}

//NewTransaction - Create new instance of transaction struct with specified arguments.
func NewTransaction(nonce uint64, to Address, amount *big.Int, data []byte, contract *contracts.Contract, extra []byte) *Transaction {
	return newTransaction(nonce, &to, amount, data, contract, extra)
}

//NewContractCreation - Create new instance of transaction struct specifying contract creation arguments.
func NewContractCreation(nonce uint64, amount *big.Int, data []byte, extra []byte) *Transaction {
	return newTransaction(nonce, nil, amount, data, nil, extra)
}

func newTransaction(nonce uint64, to *Address, amount *big.Int, data []byte, contract *contracts.Contract, extra []byte) *Transaction {
	txdata := transactiondata{
		Nonce:     nonce,
		Recipient: to,
		Payload:   data,
		Amount:    new(big.Int),
		Time:      time.Now(),
		Extra:     extra,
	}

	if amount != nil {
		txdata.Amount.Set(amount)
	}

	return &Transaction{Data: txdata, Contract: contract}
}
