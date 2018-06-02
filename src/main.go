package main

import (
	"encoding/json"
	"fmt"
	"indo-go/src/common"
	"indo-go/src/consensus"
	"indo-go/src/contracts"
	"indo-go/src/core/types"
	"indo-go/src/networking"
	"indo-go/src/networking/discovery"
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

	//Test chain serialization

	testchain.WriteChainToMemory(common.GetCurrentDir())

	testDesChain := types.ReadChainFromMemory(common.GetCurrentDir())

	//Dump deserialized chain

	b, err := json.MarshalIndent(testDesChain, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)

	//Test nodeDB serialization

	db.WriteDbToMemory(common.GetCurrentDir())

	testDb := discovery.ReadDbFromMemory(common.GetCurrentDir())

	//Dump deserialized database

	c, err := json.MarshalIndent(testDb, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(c)
}
