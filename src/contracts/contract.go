package contracts

//import "crypto/sha256"
//import "bytes"

import types "indogo/indo-go/src/core/types/payload"

//Contract - file/data representing conditions and actions filled and performed during a transaction.
type Contract struct {
	Payloads   []types.Payload `json:"payloads"`
	Identifier []byte          `json:"identifier"`
	//addressLength int
}
