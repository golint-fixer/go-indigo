package types

import (
	"reflect"

	"github.com/mitsukomegumi/indo-go/src/common"
)

// WitnessTransaction - add witness data to specified transaction if verified
func WitnessTransaction(ch *Chain, wallet *Wallet, tx *Transaction, witness *Witness) {
	if VerifyTransaction(tx) {
		tx.Weight += *CalculateWitnessWeight(witness)
		tx.Verifications++

		if tx.Verifications == 1 && string(tx.Data.Payload[:]) != "tx reward" {
			go handleReward(ch, wallet, tx, witness)
		} else if string(tx.Data.Payload[:]) == "tx reward" && !reflect.ValueOf(tx.Data.Root).IsNil() {
			*tx.Data.Root.Data.UnspentReward -= uint64(*tx.Data.Amount)
		}

		if reflect.ValueOf(tx.InitialWitness).IsNil() {
			tx.InitialWitness = witness
		}

		common.ThrowWarning("Added witness; transaction verified with weight " + common.FloatToString(tx.Weight))
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
func handleReward(ch *Chain, wallet *Wallet, tx *Transaction, witness *Witness) {
	rewardVal := float64(tx.Reward)
	nTx := NewTransaction(ch, 0, *witness.WitnessAccount, wallet.PrivateKey, wallet.PrivateKeySeeds, wallet.PublicKey, &rewardVal, []byte("tx reward"), nil, []byte("tx reward"))
	nTx.Data.Root = tx
	nTx.Reward = 0

	WitnessTransaction(ch, wallet, nTx, witness)
	(*ch).AddTransaction(nTx)

	err := Relay(nTx, ch.NodeDb)

	if err != nil {
		common.ThrowWarning(err.Error())
	}
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
		return true
	}
	return false
}
