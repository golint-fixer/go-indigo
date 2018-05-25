package types

import (
	"time"
)

//Witness - Data representation of block witness
type Witness struct {
	WitnessTime      time.Time `json:"witness timestamp"`
	WitnessedTxCount int       `json:"witness reputation"`
	WitnessSignature Signature `json:"witness signature"`
	WitnessAge       int       `json:"witness age"`
}

// getWitnessTime - return current time for use in witness
func getWitnessTime() time.Time {
	return time.Now().UTC()
}

// NewWitness - create & return new witness instance
func NewWitness(WitnessedTxCount int, WitnessSignature Signature, WitnessAge int) Witness {
	return Witness{WitnessTime: getWitnessTime(), WitnessedTxCount: WitnessedTxCount, WitnessSignature: WitnessSignature, WitnessAge: WitnessAge}
}
