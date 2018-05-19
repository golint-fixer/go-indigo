package types

import (
	"math/big"
	"time"
)

//Witness - Data representation of block witness
type Witness struct {
	WitnessTime      time.Time `json:"witness timestamp"`
	WitnessedTxCount big.Int   `json:"witness reputation"`
	WitnessSignature Signature `json:"witness signature"`
	WitnessWeight    big.Int   `json:"witness weight"`
	WitnessAge       big.Int   `json:"witness age"`
}
