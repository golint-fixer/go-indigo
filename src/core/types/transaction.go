package types

import (
	"bytes"
	"crypto"
	"encoding/json"
	"reflect"
	// crypto/sha256 - required for hashing functions
	_ "crypto/sha256"
	"fmt"
	"time"

	"github.com/mitsukomegumi/go-indigo/src/common"
	"github.com/mitsukomegumi/go-indigo/src/contracts"
)

//Transaction - Data representing transfer of value (can be null), as well as the transfer of data via payload. May be triggered on conditions, set via smart contract.
type Transaction struct {
	Data transactiondata `json:"txdata"`

	Contract *contracts.Contract `json:"contract"`

	Verifications uint64  `json:"confirmations"`
	Weight        float64 `json:"weight"`

	InitialWitness *Witness

	SendingAccount Account `json:"sending account"`

	Reward uint64 `json:"reward"`

	ChainVersion uint64 `json:"chainver"`
}

type transactiondata struct {
	// Initialized in func:
	Root          *Transaction `json:"root" gencodec:"required"`
	Nonce         uint64       `json:"nonce" gencodec:"required"`
	Recipient     *Address     `json:"recipient" gencodec:"required"`
	Amount        *float64     `json:"value" gencodec:"required"`
	UnspentReward *uint64      `json:"unspent" gencoded:"required"`
	Payload       []byte       `json:"payload" gencodec:"required"`
	Time          time.Time    `json:"timestamp" gencodec:"required"`
	Extra         []byte       `json:"extraData" gencodec:"required"`

	// Initialized at intercept:
	InitialHash *Hash `json:"hash" gencodec:"required"`
	ParentHash  *Hash `json:"parentHash" gencodec:"required"`
}

//NewTransaction - Create new instance of transaction struct with specified arguments.
func NewTransaction(ch *Chain, nonce uint64, SendingAccount Account, PrivateKey string, PrivateKeySeeds []string, to Address, amount *float64, data []byte, contract *contracts.Contract, extra []byte) *Transaction {
	if common.CheckKeys(PrivateKey, PrivateKeySeeds, SendingAccount.Address) {
		return newTransaction(ch, nonce, SendingAccount, &to, amount, data, contract, extra)
	}

	zeroVal := float64(0)

	return newTransaction(ch, nonce, SendingAccount, &to, &zeroVal, data, contract, extra)
}

//NewContractCreation - Create new instance of transaction struct specifying contract creation arguments.
func NewContractCreation(ch *Chain, nonce uint64, IssuingAccount Account, amount *float64, data []byte, extra []byte) *Transaction {
	return newTransaction(ch, nonce, IssuingAccount, nil, amount, data, nil, extra)
}

func newTransaction(Ch *Chain, nonce uint64, from Account, to *Address, amount *float64, data []byte, contract *contracts.Contract, extra []byte) *Transaction {
	hash := crypto.SHA256.New()
	txdata := transactiondata{
		Nonce:       nonce,
		Recipient:   to,
		Payload:     data,
		Amount:      new(float64),
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

	tx := Transaction{Data: txdata, Contract: contract, Weight: float64(0), Verifications: uint64(0), SendingAccount: from, Reward: 0}

	tx.Reward = tx.calculateReward(Ch)
	tx.Data.UnspentReward = &tx.Reward

	return &tx
}

// DecodeTxFromBytes - decode transaction from specified byte array, returning transaction
func DecodeTxFromBytes(wallet *Wallet, ch *Chain, wit *Witness, b []byte) *Transaction {
	plTx := Transaction{}
	err := json.NewDecoder(bytes.NewReader(b)).Decode(&plTx)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	WitnessTransaction(ch, wallet, &plTx, wit)

	return &plTx
}

func (tx Transaction) calculateReward(Ch *Chain) uint64 {
	actualCh := *Ch

	max := actualCh.MaxCirculating
	curr := actualCh.Circulating
	base := actualCh.Base

	var lastreward uint64

	if string(tx.Data.Payload[:]) == "tx reward" {
		return 0
	}

	if len(actualCh.Transactions) != 0 && !reflect.ValueOf(actualCh.Transactions[len(actualCh.Transactions)-1]).IsNil() {
		lastreward = actualCh.Transactions[len(actualCh.Transactions)-1].Reward
	} else {
		lastreward = base
	}

	common.ThrowWarning("\nlast reward: " + string(lastreward))

	if lastreward != 0 && max != 0 && base != 0 {
		if curr+(lastreward-(lastreward/200)) > max {
			return 0
		}
		return lastreward - (lastreward / 200)
	}

	if base == 0 {
		base = 10
	}

	return base
}
