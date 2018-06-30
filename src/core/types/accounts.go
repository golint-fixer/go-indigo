package types

// Account - represents account on the network
type Account struct {
	Address      Address        `json:"address"`
	Balance      int            `json:"balance"`
	URL          URL            `json:"url"`
	Transactions []*Transaction `json:"account transactions"`
}

// GetBalance - returns balance of specified account.
func GetBalance(account Account) int {
	return account.Balance
}

// NewAccount - return new account
func NewAccount(Address Address) *Account {
	return &Account{Address: Address}
}
