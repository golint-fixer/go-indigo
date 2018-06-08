package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"indo-go/src/common"
	"indo-go/src/consensus"
	"indo-go/src/contracts"
	"indo-go/src/core/types"
	"indo-go/src/networking"
	"indo-go/src/networking/discovery"
	"os"
)

var relayFlag = flag.Bool("relay", false, "used for debugging")
var listenFlag = flag.Bool("listen", false, "used for debugging")
var hostFlag = flag.Bool("host", false, "used for debugging")
var fetchFlag = flag.Bool("fetch", false, "used for debugging")
var loopFlag = flag.Bool("forever", false, "used for debugging")

func main() {
	flag.Parse()

	selfID := discovery.NodeID{} //Testing init of NodeID (self reference)

	db := discovery.NewNodeDatabase(selfID, "108.6.212.149") //Initializing net New NodeDatabase
	db.WriteDbToMemory(common.GetCurrentDir())

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

	//Test nodeDB serialization

	db.WriteDbToMemory(common.GetCurrentDir())

	fmt.Println("current dir: " + common.GetCurrentDir())

	testDb := discovery.ReadDbFromMemory(common.GetCurrentDir())

	fmt.Println("\nbest node: " + testDb.FindNode())

	fmt.Println("nodelist size: ")
	fmt.Println(len(testDb.NodeAddress))

	if *listenFlag == true {
		fmt.Println("listening")
		LatestTransaction := networking.ListenRelay()
		fmt.Println(LatestTransaction)

		// Dump fetched tx

		b, err := json.MarshalIndent(LatestTransaction, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}
		os.Stdout.Write(b)
	} else if *relayFlag == true {
		fmt.Println("attempting to relay")
		networking.Relay(test, testDb)
	} else if *hostFlag == true {
		fmt.Println("attempting to host")
		networking.HostChain(testDesChain, testDb, *loopFlag)
	} else if *fetchFlag == true {
		fmt.Println("attempting to fetch chain")
		networking.FetchChainWithAdd(testDesChain, testDb)

		// Dump fetched chain

		b, err := json.MarshalIndent(testDesChain, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}
		os.Stdout.Write(b)
	}
}
