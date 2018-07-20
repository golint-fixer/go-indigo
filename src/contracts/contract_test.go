package contracts

import (
	"strconv"
	"testing"
)

func TestContract(t *testing.T) {
	in := 1
	out := 2

	var1 := ContractVariable{Input: float64(in), Output: float64(out), Conditionals: []byte("=="), Modifier: float64(1), ModifierOperation: []byte("+")}

	environment := ContractEnvironment{Variables: []*ContractVariable{&var1}}

	contract := Contract{RuntimeEnv: environment}

	checkedCondition := contract.RuntimeEnv.Variables[0].CheckCondition()

	if checkedCondition == false {
		t.Errorf("Contract invalid: %s", strconv.FormatBool(checkedCondition))
	} else {
		t.Logf("Contract valid: %s", strconv.FormatBool(checkedCondition))
	}
}
