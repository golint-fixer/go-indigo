package types

import (
	"reflect"
	"strings"

	"github.com/mitsukomegumi/indo-go/src/common"
)

// WitnessTransaction - add witness data to specified transaction if verified
func WitnessTransaction(ch *Chain, wallet *Wallet, tx *Transaction, witness *Witness) {
	if VerifyTransaction(tx) {
		tx.Weight += *CalculateWitnessWeight(witness)
		tx.Verifications++

		if tx.Verifications == 2 && string(tx.Data.Payload[:]) != "tx reward" {
			common.ThrowSuccess("transaction verified; creating reward")
			go handleReward(ch, wallet, tx, witness)
		} else if string(tx.Data.Payload[:]) == "tx reward" && !reflect.ValueOf(tx.Data.Root).IsNil() {
			*tx.Data.Root.Data.UnspentReward -= uint64(*tx.Data.Amount)
		}

		if reflect.ValueOf(tx.InitialWitness).IsNil() {
			tx.InitialWitness = witness
		}

		if string(tx.Data.Payload[:]) == "tx reward" {
			common.ThrowSuccess("Added witness; reward verified with weight " + common.FloatToString(tx.Weight))
		} else {
			common.ThrowWarning("Added witness; transaction verified with weight " + common.FloatToString(tx.Weight))
		}
	} else {
		tx.Weight -= *CalculateWitnessWeight(witness)
		tx.Verifications++

		if reflect.ValueOf(tx.InitialWitness).IsNil() {
			tx.InitialWitness = witness
		}

		common.ThrowWarning("Added witness, removed weight; transaction illegitimate with weight " + common.FloatToString(tx.Weight))
	}
}

// handleReward - perform needed actions to account for a reward
func handleReward(ch *Chain, wallet *Wallet, tx *Transaction, witness *Witness) error {
	go func() error {
		fetchChain, err := FetchChain(ch.NodeDb)

		if err != nil && !strings.Contains(err.Error(), "an existing connection") {
			return err
		}

		(*ch) = *fetchChain

		return nil
	}()

	rewardVal := float64(tx.Reward)
	nTx := NewTransaction(ch, 0, *witness.WitnessAccount, wallet.PrivateKey, wallet.PrivateKeySeeds, wallet.PublicKey, &rewardVal, []byte("tx reward"), nil, []byte("tx reward"))
	nTx.Data.Root = tx
	nTx.Reward = 0

	WitnessTransaction(ch, wallet, nTx, witness)
	(*ch).AddTransaction(nTx)
	(*ch).Circulating += nTx.Data.Root.Reward

	err := Relay(nTx, ch.NodeDb)

	if err != nil && strings.Contains(err.Error(), "; fetch") {
		return err
	}

	return nil
}

// CalculateWitnessWeight - calculate weight for individual witness based on implied or given weight
func CalculateWitnessWeight(witness *Witness) *float64 {
	witnessWeight := float64(witness.WitnessedTxCount / witness.WitnessAge)
	return &witnessWeight
}

// VerifyTransaction - checks validity of transaction, returning bool
func VerifyTransaction(tx *Transaction) bool {
	balance := GetBalance(tx.SendingAccount)
	amountTransacted := *tx.Data.Amount

	if balance <= amountTransacted {
		return true
	} else if string(tx.Data.Payload[:]) == "tx reward" {
		return verifyCoinbase(tx)
	}
	return false
}

func verifyCoinbase(tx *Transaction) bool {
	if float64(*tx.Data.Root.Data.UnspentReward) == *tx.Data.Amount && string(tx.Data.Payload[:]) == "tx reward" {
		common.ThrowSuccess("reward coinbase verified")
		return true
	}
	return false
}
