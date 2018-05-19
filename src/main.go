package main

import (
	"encoding/json"
	"fmt"
	"indogo/src/contracts"
	types "indogo/src/core/types"
	"math/big"
	"os"
)

func main() {
	testcontract := new(contracts.Contract)
	test := types.NewTransaction(uint64(1), types.HexToAddress("281055afc982d96fab65b3a49cac8b878184cb16"), big.NewInt(1000000), []byte{0x11, 0x11, 0x11}, testcontract, nil)

	b, err := json.MarshalIndent(test, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
}
