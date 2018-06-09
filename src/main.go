package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitsukomegumi/indo-go/src/common"
	"github.com/mitsukomegumi/indo-go/src/consensus"
	"github.com/mitsukomegumi/indo-go/src/contracts"
	"github.com/mitsukomegumi/indo-go/src/core/types"
	"github.com/mitsukomegumi/indo-go/src/networking"
	"github.com/mitsukomegumi/indo-go/src/networking/discovery"
)

var relayFlag = flag.Bool("relay", false, "used for debugging")
var listenFlag = flag.Bool("listen", false, "used for debugging")
var hostFlag = flag.Bool("host", false, "used for debugging")
var fetchFlag = flag.Bool("fetch", false, "used for debugging")
var loopFlag = flag.Bool("forever", false, "used for debugging")

func main() {
	flag.Parse()

	if *listenFlag || *hostFlag {

		eDb := discovery.NodeDatabase{}
		rErr := common.ReadGob(common.GetCurrentDir()+"nodeDb.gob", eDb)

		if strings.Contains(rErr.Error(), "cannot find") {
			common.ThrowWarning(rErr.Error())

			gd, err := networking.GetGateway()
			ip, err := gd.ExternalIP()

			if err != nil {
				panic(err)
			}

			networking.PrepareForConnection(gd)

			selfID := discovery.NodeID{} //Testing init of NodeID (self reference)

			db := discovery.NewNodeDatabase(selfID, ip) //Initializing net New NodeDatabase
			db.WriteDbToMemory(common.GetCurrentDir())
		}
	} else {
		selfID := discovery.NodeID{} //Testing init of NodeID (self reference)

		db := discovery.NewNodeDatabase(selfID, "") //Initializing net New NodeDatabase
		db.WriteDbToMemory(common.GetCurrentDir())
	}

	db := discovery.ReadDbFromMemory(common.GetCurrentDir())

	if *relayFlag || *hostFlag {
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

		if *relayFlag == true {
			fmt.Println("attempting to relay")
			networking.Relay(test, db)
		} else if *hostFlag == true {
			fmt.Println("attempting to host")
			networking.HostChain(testDesChain, db, *loopFlag)
		}
	}

	fmt.Println("\nbest node: " + db.FindNode())

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
	} else if *fetchFlag == true {
		testDesChain := types.Chain{}
		fmt.Println("attempting to fetch chain")
		networking.FetchChainWithAdd(&testDesChain, db)

		// Dump fetched chain

		b, err := json.MarshalIndent(testDesChain, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}
		os.Stdout.Write(b)
	}
}
