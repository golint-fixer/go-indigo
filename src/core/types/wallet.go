package types

import (
	"errors"

	"github.com/mitsukomegumi/go-indigo/src/common"
)

// Wallet - holds private, public keys linking to specified account
type Wallet struct {
	PrivateKeySeeds []string `json:"privatekeyseeds"`
	PrivateKey      string   `json:"privatekeys"`
	PublicKey       Address  `json:"publickeys"` // Also known as a wallet address

	Balance float64 `json:"balance"`

	OriginVersion uint64 `json:"origin"`  // Chain version at wallet creation
	LastVersion   uint64 `json:"version"` // Last scanned block with version number

	Transactions []*Transaction `json:"transactions"`

	Account *Account `json:"account"`
}

// ClaimWallet - verifies private keys of specified public key, returns wallet of specified keys
func ClaimWallet(Ch *Chain, pub Address, Private string, PrivateKeySeeds []string) (Wallet, error) {
	if common.CheckKeys(Private, PrivateKeySeeds, pub) {
		wallet := Wallet{PrivateKeySeeds: PrivateKeySeeds, PrivateKey: Private, PublicKey: pub, Balance: 0, Account: NewAccount(pub)}
		wallet.ScanChain(Ch)
		return wallet, nil
	}
	return Wallet{}, errors.New("invalid keys")
}

// NewWallet - create new wallet instance
func NewWallet(Ch *Chain) *Wallet {
	tempWallet := new(Wallet)

	seed := tempWallet.generateSeed()
	tempWallet.PrivateKeySeeds = seed

	private := tempWallet.generatePrivateKey()
	tempWallet.PrivateKey = private

	public := tempWallet.generatePublicKey()
	tempWallet.PublicKey = BytesToAddress([]byte(public))

	wallet := Wallet{PrivateKeySeeds: tempWallet.generateSeed(), PrivateKey: tempWallet.generatePrivateKey(), PublicKey: BytesToAddress([]byte(tempWallet.generatePublicKey())), Balance: 0, OriginVersion: Ch.Version, LastVersion: Ch.Version}

	wallet.Account = NewAccount(wallet.PublicKey)

	return &wallet
}

// ScanChain - search specified chain for transactions with public key
func (wallet Wallet) ScanChain(Ch *Chain) {
	wallet.findSent(Ch)
	wallet.findReceived(Ch)

	wallet.Account.Balance = wallet.Balance
}

func (wallet Wallet) findSent(Ch *Chain) {
	x := wallet.LastVersion

	for x != uint64(len(Ch.Transactions)) {
		if Ch.Transactions[x].SendingAccount.Address == wallet.PublicKey && string(Ch.Transactions[x].Data.Payload[:]) != "tx reward" {
			wallet.Transactions = append(wallet.Transactions, Ch.Transactions[x])
			wallet.Balance -= *Ch.Transactions[x].Data.Amount
		}
		x++
	}

	wallet.LastVersion = x
}

func (wallet Wallet) findReceived(Ch *Chain) {
	x := wallet.LastVersion

	for x != uint64(len(Ch.Transactions)) {
		if *Ch.Transactions[x].Data.Recipient == wallet.PublicKey {
			wallet.Transactions = append(wallet.Transactions, Ch.Transactions[x])
			wallet.Balance += *Ch.Transactions[x].Data.Amount
		}
		x++
	}

	wallet.LastVersion = x
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

// WriteToMemory - Write specified wallet to memory
func (wallet Wallet) WriteToMemory(dir string) error {
	err := common.WriteGob(dir, wallet)

	if err != nil {
		return err
	}

	return nil
}

// ReadWalletFromMemory - read wallet from specified directory
func ReadWalletFromMemory(dir string) (Wallet, error) {
	wallet := Wallet{}
	err := common.ReadGob(dir, wallet)

	if err != nil {
		return wallet, err
	}

	return wallet, nil
}
