package types

import (
	"github.com/mitsukomegumi/indo-go/src/common"
)

// Account - represents account on the network
type Account struct {
	Address      common.Address `json:"address"`
	URL          URL            `json:"url"`
	Transactions []*Transaction `json:"account transactions"`
}

// GetBalance - returns balance of specified account.
func GetBalance(account Account) *int {
	test := 0
	return &test
}

// NewAccount - return new account
func NewAccount(Address common.Address) *Account {
	return &Account{Address: Address}
}
