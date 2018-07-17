package contracts

import (
	"fmt"
	"testing"
)

func TestContract(t *testing.T) {
	in := 1
	out := 1

	var1 := ContractVariable{Input: in, Output: out, Conditionals: []byte("==")}

	environment := ContractEnvironment{Variables: []*ContractVariable{&var1}}

	contract := Contract{RuntimeEnv: environment}

	fmt.Println(contract.RuntimeEnv.Variables[0].CheckCondition())
}
