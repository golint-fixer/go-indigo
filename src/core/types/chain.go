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

	Version int `json:"version"`
}

// AddTransaction - Add transaction to specified chain object
func (RefChain *Chain) AddTransaction(Transaction *Transaction) {
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
func (RefChain Chain) WriteChainToMemory(path string) {
	common.WriteGob(path+string(RefChain.Identifier)+"Chain.gob", RefChain)
	RefChain.NodeDb.WriteDbToMemory(common.GetCurrentDir())
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
func DecodeChainFromBytes(b []byte) *Chain {
	plCh := Chain{}
	err := json.NewDecoder(bytes.NewReader(b)).Decode(&plCh)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return &plCh
}
