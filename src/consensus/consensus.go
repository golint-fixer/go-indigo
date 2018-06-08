package consensus

import (
	"reflect"

	"github.com/mitsukomegumi/indo-go/src/common"
	"github.com/mitsukomegumi/indo-go/src/core/types"
)

// WitnessTransaction - add witness data to specified transaction if verified
func WitnessTransaction(tx *types.Transaction, witness *types.Witness) {
	if VerifyTransaction(tx) {
		tx.Weight += *CalculateWitnessWeight(witness)
		tx.Verifications++

		if reflect.ValueOf(tx.InitialWitness).IsNil() {
			tx.InitialWitness = witness
		}

		common.ThrowWarning("Added witness; transaction verified with weight " + string(tx.Weight))
	} else {
		tx.Weight -= *CalculateWitnessWeight(witness)
		tx.Verifications++

		if reflect.ValueOf(tx.InitialWitness).IsNil() {
			tx.InitialWitness = witness
		}

		common.ThrowWarning("Added witness, removed weight; transaction illegitimate with weight " + string(tx.Weight))
	}
}

// CalculateWeight - calculate weight for transaction based on current weight or implied weight
func CalculateWeight(tx *types.Transaction) {

}

// CalculateWitnessWeight - calculate weight for individual witness based on implied or given weight
func CalculateWitnessWeight(witness *types.Witness) *int {
	witnessWeight := int(witness.WitnessedTxCount / witness.WitnessAge)
	return &witnessWeight
}

// VerifyTransaction - checks validity of transaction, returning bool
func VerifyTransaction(tx *types.Transaction) bool {
	balance := *types.GetBalance(tx.SendingAccount)
	amountTransacted := *tx.Data.Amount

	if balance <= amountTransacted {
		return true
	}
	return false
}
