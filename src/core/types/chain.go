package types

import (
	"fmt"
	"indogo/src/contracts"
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
