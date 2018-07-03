package types

import (
	"bytes"
	"crypto"
	"encoding/json"
	// crypto/sha256 - required for hashing functions
	_ "crypto/sha256"
	"fmt"
	"time"

	"github.com/mitsukomegumi/indo-go/src/common"
	"github.com/mitsukomegumi/indo-go/src/contracts"
)

//Transaction - Data representing transfer of value (can be null), as well as the transfer of data via payload. May be triggered on conditions, set via smart contract.
type Transaction struct {
	Data transactiondata `json:"txdata"`

	Contract *contracts.Contract `json:"contract"`

	Verifications int `json:"confirmations"`
	Weight        int `json:"weight"`

	InitialWitness *Witness

	SendingAccount Account `json:"sending account"`

	Reward int `json:"reward"`

	ChainVersion int `json:"chainver"`
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
func NewTransaction(ch *Chain, nonce uint64, SendingAccount Account, PrivateKey string, PrivateKeySeeds []string, to Address, amount *int, data []byte, contract *contracts.Contract, extra []byte) *Transaction {
	if common.CheckKeys(PrivateKey, PrivateKeySeeds, string(SendingAccount.Address[:])) {
		return newTransaction(ch, nonce, SendingAccount, &to, amount, data, contract, extra)
	}

	zeroVal := 0

	return newTransaction(ch, nonce, SendingAccount, &to, &zeroVal, data, contract, extra)
}

//NewContractCreation - Create new instance of transaction struct specifying contract creation arguments.
func NewContractCreation(ch *Chain, nonce uint64, IssuingAccount Account, amount *int, data []byte, extra []byte) *Transaction {
	return newTransaction(ch, nonce, IssuingAccount, nil, amount, data, nil, extra)
}

func newTransaction(Ch *Chain, nonce uint64, from Account, to *Address, amount *int, data []byte, contract *contracts.Contract, extra []byte) *Transaction {
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

	tx := Transaction{Data: txdata, Contract: contract, Weight: int(0), Verifications: int(0), SendingAccount: from, Reward: 0}

	tx.Reward = tx.calculateReward(Ch)

	(*Ch).Circulating += tx.Reward

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

func (tx Transaction) calculateReward(Ch *Chain) int {
	max := Ch.MaxCirculating
	curr := Ch.Circulating
	base := Ch.Base
	lastreward := Ch.Transactions[len(Ch.Transactions)-1].Reward

	if lastreward != 0 && max != 0 && base != 0 {
		if curr+lastreward/2 > max {
			return 0
		}
		return lastreward / 2
	}

	base = 10
	(*Ch).Base = 10

	return base
}
