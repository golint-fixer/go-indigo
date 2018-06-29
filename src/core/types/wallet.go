package types

import (
	"github.com/mitsukomegumi/indo-go/src/common"
)

// Wallet - holds private, public keys linking to specified account
type Wallet struct {
	PrivateKeySeeds []string `json:"privatekeyseeds"`
	PrivateKey      string   `json:"privatekeys"`
	PublicKey       Address  `json:"publickeys"` // Also known as a wallet address

	Balance int `json:"balance"`

	OriginVersion int `json:"origin"`  // Chain version at wallet creation
	LastVersion   int `json:"version"` // Last scanned block with version number

	Transactions []*Transaction `json:"transactions"`
}

// NewWallet - create new wallet instance
func NewWallet(Ch *Chain) *Wallet {
	tempWallet := new(Wallet)
	wallet := Wallet{PrivateKeySeeds: tempWallet.generateSeed(), PrivateKey: tempWallet.generatePrivateKey(), PublicKey: BytesToAddress([]byte(tempWallet.generatePublicKey())), Balance: 0, OriginVersion: Ch.Version, LastVersion: Ch.Version}

	return &wallet
}

// ScanChain - search specified chain for transactions with public key
func (wallet Wallet) ScanChain(Ch *Chain) {
	wallet.findSent(Ch)
	wallet.findReceived(Ch)
}

func (wallet Wallet) findSent(Ch *Chain) []*Transaction {
	x := wallet.LastVersion

	for x != len(Ch.Transactions) {
		if Ch.Transactions[x].SendingAccount.Address == wallet.PublicKey {
			wallet.Transactions = append(wallet.Transactions, Ch.Transactions[x])
		}
		x++
	}

	return nil
}

func (wallet Wallet) findReceived(Ch *Chain) []*Transaction {
	x := wallet.LastVersion

	for x != len(Ch.Transactions) {
		if *Ch.Transactions[x].Data.Recipient == wallet.PublicKey {
			wallet.Transactions = append(wallet.Transactions, Ch.Transactions[x])
		}
		x++
	}

	return nil
}

func (wallet Wallet) generatePublicKey() string {
	combined := wallet.PrivateKeySeeds[0] + wallet.PrivateKeySeeds[1]
	return common.SHA256([]byte(wallet.PrivateKey + combined))
}

func (wallet Wallet) generatePrivateKey() string {
	combined := wallet.PrivateKeySeeds[0] + wallet.PrivateKeySeeds[1]
	return common.SHA256([]byte(combined))
}

func (wallet Wallet) generateSeed() []string {
	timeStamp := common.GetCurrentTime()
	randStr := common.RandStringRunes(8)

	return []string{timeStamp.String(), randStr}
}
