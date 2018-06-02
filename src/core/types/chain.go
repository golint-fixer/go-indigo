package types

import (
	"fmt"
	"indo-go/src/common"
	"indo-go/src/contracts"
	"reflect"
)

// Chain - Connected collection of transactions
type Chain struct {
	ParentContract *contracts.Contract `json:"parentcontract"`
	Identifier     Identifier

	Transactions []*Transaction
}

// AddTransaction - Add transaction to specified chain object
func (RefChain *Chain) AddTransaction(Transaction *Transaction) {
	RefChain.Transactions = append(RefChain.Transactions, Transaction)
	fmt.Println("Transaction added to chain")
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
