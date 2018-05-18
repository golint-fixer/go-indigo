package main

import (
	"fmt"
	"indogo/indo-go/src/contracts"
	types "indogo/indo-go/src/core/types"
	"log"
	"math/big"
	"strings"
)

func main() {
	testcontract := new(contracts.Contract)
	test := types.NewTransaction(uint64(1), types.BytesToAddress([]byte{0x11}), big.NewInt(1000000), []byte{0x11, 0x11, 0x11}, testcontract, nil)
	log.Println(strings.Replace(fmt.Sprintf("%#v", test), ", ", "\n", -1))
}
