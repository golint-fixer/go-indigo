package witness

import (
	"math/big"
	"time"
)

//Witness - Data representation of block witness
type Witness struct {
	Witnesstime      time.Time
	WitnessedTxCount big.Int
}
