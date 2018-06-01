package main

import (
	"encoding/json"
	"fmt"
	"indogo/src/common"
	"indogo/src/consensus"
	"indogo/src/contracts"
	"indogo/src/core/types"
	"indogo/src/networking"
	"indogo/src/networking/discovery"
	"os"
)

func main() {
	selfID := networking.NodeID{}

	db := discovery.NewNodeDatabase(selfID)
	db.AddNode("1.1.1.1", selfID)

	accountAddress := common.HexToAddress("281055afc982d96fab65b3a49cac8b878184cb16")

	account := types.NewAccount(accountAddress)

	signature := types.HexToSignature("281055afc982d96fab65b3a49cac8b878184cb16")

	witness := types.NewWitness(1000, signature, 100)

	testcontract := new(contracts.Contract)

	testchain := types.Chain{ParentContract: testcontract}

	test := types.NewTransaction(uint64(1), *account, types.HexToAddress("281055afc982d96fab65b3a49cac8b878184cb16"), common.IntToPointer(1000), []byte{0x11, 0x11, 0x11}, testcontract, nil)

	consensus.WitnessTransaction(test, &witness)

	testchain.AddTransaction(test)

	b, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)

	db.WriteDbToMemory("R:\\gocode\\src\\indogo\\src\\globbityglob.gob")

	testDb := new(discovery.NodeDatabase)

	error := common.ReadGob("R:\\gocode\\src\\indogo\\src\\globbityglob.gob", testDb)
	if error != nil {
		fmt.Println(error)
	}

	c, err := json.MarshalIndent(testDb, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(c)
}
