package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mitsukomegumi/indo-go/src/common"
	"github.com/mitsukomegumi/indo-go/src/contracts"
	"github.com/mitsukomegumi/indo-go/src/networking/discovery"
)

// Chain - Connected collection of transactions
type Chain struct {
	ParentContract *contracts.Contract `json:"parentcontract"`
	Identifier     Identifier          `json:"identifier"`

	NodeDb *discovery.NodeDatabase `json:"database"`

	Transactions []*Transaction `json:"transactions"`

	MaxCirculating uint64 `json:"maxcirculating"`
	Circulating    uint64 `json:"circulating"`

	Base uint64 `json:"base"`

	Version uint64 `json:"version"`
}

// AddTransaction - Add transaction to specified chain object
func (RefChain *Chain) AddTransaction(Transaction *Transaction) {
	if VerifyTransaction(Transaction) {
		RefChain.Transactions = append(RefChain.Transactions, Transaction)

		if Transaction.ChainVersion == 0 {
			RefChain.Version++
			Transaction.ChainVersion = RefChain.Version
		} else if Transaction.ChainVersion > RefChain.Version {
			if (Transaction.ChainVersion - RefChain.Version) == 1 {
				RefChain.Version = Transaction.ChainVersion
			}
		}

		fmt.Println("transaction added to chain")

		RefChain.WriteChainToMemory(common.GetCurrentDir())
	}
}

// FindUnverifiedTransactions - Browse chain for most recent unverified transactions
func (RefChain Chain) FindUnverifiedTransactions(TxCount int) []*Transaction {

	var UnverifiedTransactions []*Transaction

	x := len(RefChain.Transactions) - 1

	targetTxCount := TxCount + 1

	for len(UnverifiedTransactions) != targetTxCount {
		if reflect.ValueOf(RefChain.Transactions[x].InitialWitness).IsNil() || reflect.ValueOf(RefChain.Transactions[x].InitialWitness).IsNil() {
			UnverifiedTransactions = append(UnverifiedTransactions, RefChain.Transactions[x])
		} else {
			x--
		}
	}

	return UnverifiedTransactions
}

// WriteChainToMemory - create serialized instance of specified chain in specified path (string)
func (RefChain Chain) WriteChainToMemory(path string) error {
	common.WriteGob(path+string(RefChain.Identifier)+"Chain.gob", RefChain)
	err := RefChain.NodeDb.WriteDbToMemory(common.GetCurrentDir())

	if err != nil {
		return err
	}
	return nil
}

// ReadChainFromMemory - read serialized object of specified chain from specified path
func ReadChainFromMemory(path string) *Chain {
	tempChain := new(Chain)

	error := common.ReadGob(path+"Chain.gob", tempChain)
	if error != nil {
		fmt.Println(error)
	} else {
		return tempChain
	}
	return nil
}

// DecodeChainFromBytes - decode chain from specified byte array, returning new chain
func DecodeChainFromBytes(b []byte) (*Chain, error) {
	plCh := Chain{}
	err := json.NewDecoder(bytes.NewReader(b)).Decode(&plCh)

	if err != nil {
		return nil, err
	}

	return &plCh, nil
}
