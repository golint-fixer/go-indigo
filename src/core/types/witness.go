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
