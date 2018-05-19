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
	test := types.NewTransaction(uint64(1), types.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"), big.NewInt(1000000), []byte{0x11, 0x11, 0x11}, testcontract, nil)

	b, err := json.MarshalIndent(test, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
}
