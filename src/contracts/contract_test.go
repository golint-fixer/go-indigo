package contracts

import (
	"strconv"
	"testing"
)

func TestContract(t *testing.T) {
	in := 1
	out := 1

	var1 := ContractVariable{Input: in, Output: out, Conditionals: []byte("==")}

	environment := ContractEnvironment{Variables: []*ContractVariable{&var1}}

	contract := Contract{RuntimeEnv: environment}

	t.Log("contract value: " + strconv.FormatBool(contract.RuntimeEnv.Variables[0].CheckCondition()))
}
