package types

import (
	"bytes"
	"crypto"
	"encoding/json"
	// crypto/sha256 - required for hashing functions
	_ "crypto/sha256"
	"fmt"
	"sync/atomic"
	"time"

	contracts "github.com/mitsukomegumi/indo-go/src/contracts"
)

//Transaction - Data representing transfer of value (can be null), as well as the transfer of data via payload. May be triggered on conditions, set via smart contract.
type Transaction struct {
	Data transactiondata `json:"txdata"`

	Contract *contracts.Contract `json:"contract"`

	Verifications int `json:"confirmations"`
	Weight        int `json:"weight"`

	InitialWitness *Witness

	SendingAccount Account `json:"sending account"`

	ChainVersion int `json:"chainver"`

	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

type transactiondata struct {
	// Initialized in func:
	Nonce     uint64    `json:"nonce" gencodec:"required"`
	Recipient *Address  `json:"recipient"`
	Amount    *int      `json:"value" gencodec:"required"`
	Payload   []byte    `json:"payload" gencodec:"required"`
	Time      time.Time `json:"timestamp" gencodec:"required"`
	Extra     []byte    `json:"extraData" gencodec:"required"`

	// Initialized at intercept:
	InitialHash *Hash `json:"hash" gencodec:"required"`
	ParentHash  *Hash `json:"parentHash" gencodec:"required"`
}

//NewTransaction - Create new instance of transaction struct with specified arguments.
func NewTransaction(nonce uint64, SendingAccount Account, to Address, amount *int, data []byte, contract *contracts.Contract, extra []byte) *Transaction {
	return newTransaction(nonce, SendingAccount, &to, amount, data, contract, extra)
}

//NewContractCreation - Create new instance of transaction struct specifying contract creation arguments.
func NewContractCreation(nonce uint64, IssuingAccount Account, amount *int, data []byte, extra []byte) *Transaction {
	return newTransaction(nonce, IssuingAccount, nil, amount, data, nil, extra)
}

func newTransaction(nonce uint64, from Account, to *Address, amount *int, data []byte, contract *contracts.Contract, extra []byte) *Transaction {
	hash := crypto.SHA256.New()
	txdata := transactiondata{
		Nonce:       nonce,
		Recipient:   to,
		Payload:     data,
		Amount:      new(int),
		Time:        time.Now().UTC(),
		Extra:       extra,
		InitialHash: &Hash{},
	}

	s := fmt.Sprintf("%v", txdata)
	bArray := hash.Sum([]byte(s))

	*txdata.InitialHash = BytesToHash(bArray)

	if amount != nil {
		txdata.Amount = amount
	}

	return &Transaction{Data: txdata, Contract: contract, Weight: int(0), Verifications: int(0), SendingAccount: from}
}

// DecodeTxFromBytes - decode transaction from specified byte array, returning transaction
func DecodeTxFromBytes(b []byte) *Transaction {
	plTx := Transaction{}
	err := json.NewDecoder(bytes.NewReader(b)).Decode(&plTx)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return &plTx
}
