package types

import (
	"indogo/src/contracts"
)

// Chain - Connected collection of transactions
type Chain struct {
	ParentContract *contracts.Contract `json:"parentcontract"`
	Identifier     Identifier

	Transactions []*Transaction
}

func AddTransaction(Transaction *Transaction) {

}
