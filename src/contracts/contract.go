package contracts

//import "crypto/sha256"
//import "bytes"

import (
	"strings"

	types "github.com/mitsukomegumi/go-indigo/src/core/types/payload"
)

// EventTypes - Available event types
var EventTypes = []string{"met", "waiting", "failed", "invalid"}

//Contract - file/data representing conditions and actions filled and performed during a transaction.
type Contract struct {
	Payloads   []types.Payload `json:"payloads"` // Miscellaneous metadata associated with contract
	Identifier []byte          `json:"identifier"`

	RuntimeEnv ContractEnvironment `json:"environment"`
}

// ContractEnvironment - holds metadata, technical specifications of contract's runtime environment
type ContractEnvironment struct {
	Variables   []*ContractVariable `json:"contract variables"`
	MemoryAddrs []string            `json:"addresses"`
}

// ContractVariable - individual set of inputs, outputs specifying conditionals to be checked
type ContractVariable struct {
	Identifier []byte `json:"identifier"`

	Input  interface{} `json:"input"`
	Output interface{} `json:"output"`

	Conditionals []byte `json:"conditionals"`
	Modifiers    []byte `json:"modifiers"`

	Events []*ContractEvent `json:"variable events"`

	MemAddr string `json:"address"`
}

// ContractEvent - represents an update to a contrat variable's status
type ContractEvent struct {
	Identifier []byte `json:"identifier"`

	EventType string `json:"event type"`

	Action []byte `json:"action"`
}

// CheckCondition - checks whether value of specified variable is true
func (variable *ContractVariable) CheckCondition() bool {
	if strings.Contains(string(variable.Conditionals[:]), "==") {
		return variable.equalEqual()
	} else if strings.Contains(string(variable.Conditionals[:]), ">=") {
		return variable.greaterEqual()
	} else if strings.Contains(string(variable.Conditionals[:]), ">") {
		return variable.greater()
	} else if strings.Contains(string(variable.Conditionals[:]), "<=") {
		return variable.lessthanEqual()
	} else if strings.Contains(string(variable.Conditionals[:]), "<") {
		return variable.lessthan()
	}

	return false
}
