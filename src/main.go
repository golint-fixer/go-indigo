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
	selfID := networking.NodeID{} //Testing init of NodeID (self reference)

	db := discovery.NewNodeDatabase(selfID) //Initializing net New NodeDatabase
	db.AddNode("1.1.1.1", selfID)           //Adding node to database

	//Creating new account:

	accountAddress := common.HexToAddress("281055afc982d96fab65b3a49cac8b878184cb16")
	account := types.NewAccount(accountAddress)

	//Creating witness data:

	signature := types.HexToSignature("281055afc982d96fab65b3a49cac8b878184cb16")
	witness := types.NewWitness(1000, signature, 100)

	//Creating transaction, contract, chain

	testcontract := new(contracts.Contract)
	testchain := types.Chain{ParentContract: testcontract}
	test := types.NewTransaction(uint64(1), *account, types.HexToAddress("281055afc982d96fab65b3a49cac8b878184cb16"), common.IntToPointer(1000), []byte{0x11, 0x11, 0x11}, testcontract, nil)

	//Adding witness, transaction to chain

	consensus.WitnessTransaction(test, &witness)
	testchain.AddTransaction(test)

	//Dump chain

	b, err := json.MarshalIndent(testchain, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)

	//Test nodeDB serialization methods

	db.WriteDbToMemory("R:\\gocode\\src\\indogo\\src\\")

	testDb := discovery.ReadDbFromMemory("R:\\gocode\\src\\indogo\\src\\")

	//Dump read database

	c, err := json.MarshalIndent(testDb, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(c)
}
