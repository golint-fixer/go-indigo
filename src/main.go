package main

import (
	"encoding/json"
	"fmt"
	"indogo/src/common"
	"indogo/src/consensus"
	"indogo/src/contracts"
	types "indogo/src/core/types"
	"os"
)

func main() {
	accountAddress := common.HexToAddress("281055afc982d96fab65b3a49cac8b878184cb16")

	account := types.NewAccount(accountAddress)

	signature := types.HexToSignature("281055afc982d96fab65b3a49cac8b878184cb16")

	witness := types.NewWitness(1000, signature, 100)

	testcontract := new(contracts.Contract)
	test := types.NewTransaction(uint64(1), *account, types.HexToAddress("281055afc982d96fab65b3a49cac8b878184cb16"), common.IntToPointer(1000), []byte{0x11, 0x11, 0x11}, testcontract, nil)

	consensus.WitnessTransaction(test, &witness)

	b, err := json.MarshalIndent(test, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
}
