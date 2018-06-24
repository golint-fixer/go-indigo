package main

import (
	"bufio"
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

var relayFlag = flag.Bool("relay", false, "relay tx to node")
var listenFlag = flag.Bool("listen", false, "listen for transaction relays")
var hostFlag = flag.Bool("host", false, "host current copy of chain")
var fetchFlag = flag.Bool("fetch", false, "fetch current copy of chain")
var newChainFlag = flag.Bool("new", false, "create new chain")
var loopFlag = flag.Bool("forever", false, "perform indefinitely")
var fullChainFlag = flag.Bool("relaychain", false, "relay entire chain")
var registerNode = flag.Bool("regnode", false, "registers node")
var noUpNPFlag = flag.Bool("noupnp", false, "used for nodes without upnp")

/*
	TODO:
		[DONE] - test node db serialization
		[DONE] - add version to chain struct (increments on each transaction)
		[DONE] - add unit testing
*/

func main() {
	flag.Parse()

	if *relayFlag || *listenFlag || *hostFlag || *fetchFlag || *loopFlag || *fullChainFlag || *noUpNPFlag {
		if *listenFlag || *hostFlag {
			gd, err := networking.GetGateway()
			tsfRef := discovery.NodeID{}
			eDb, err := discovery.NewNodeDatabase(tsfRef, "")

			if err != nil {
				panic(err)
			}

			rErr := common.ReadGob(common.GetCurrentDir()+"nodeDb.gob", &eDb)

			if err != nil {
				panic(err)
			}

			if rErr != nil && strings.Contains(rErr.Error(), "cannot find") {
				common.ThrowWarning(rErr.Error())

				ip, err := gd.ExternalIP()
				if err != nil {
					panic(err)
				}

				selfID := discovery.NodeID{} //Testing init of NodeID (self reference)

				db, err := discovery.NewNodeDatabase(selfID, ip) //Initializing net New NodeDatabase

				if err != nil {
					panic(err)
				}

				db.WriteDbToMemory(common.GetCurrentDir())
			}
			if !*noUpNPFlag {
				fmt.Println("configuring upnp devices")
				networking.PrepareForConnection(gd, eDb)
			}
		} else {
			selfID := discovery.NodeID{} //Testing init of NodeID (self reference)

			db, err := discovery.NewNodeDatabase(selfID, "") //Initializing net New NodeDatabase

			if err != nil {
				panic(err)
			}

			db.WriteDbToMemory(common.GetCurrentDir())
		}

		db, err := discovery.ReadDbFromMemory(common.GetCurrentDir())

		if err != nil {
			panic(err)
		}

		fmt.Println("\nbest node: " + db.FindNode())

		if *relayFlag || *hostFlag || *fullChainFlag {
			//Creating new account:

			accountAddress := common.HexToAddress("281055afc982d96fab65b3a49cac8b878184cb16")
			account := types.NewAccount(accountAddress)

			//Creating witness data:

			signature := types.HexToSignature("281055afc982d96fab65b3a49cac8b878184cb16")
			witness := types.NewWitness(1000, signature, 100)

			//Creating transaction, contract, chain

			testchain := types.ReadChainFromMemory(common.GetCurrentDir())

			test := types.NewTransaction(uint64(1), *account, types.HexToAddress("281055afc982d96fab65b3a49cac8b878184cb16"), common.IntToPointer(1000), []byte{0x11, 0x11, 0x11}, nil, nil)

			//Adding witness, transaction to chain

			consensus.WitnessTransaction(test, &witness)
			testchain.AddTransaction(test)

			//Test chain serialization

			testchain.WriteChainToMemory(common.GetCurrentDir())

			testDesChain := types.ReadChainFromMemory(common.GetCurrentDir())

			if *relayFlag {
				fmt.Println("attempting to relay")
				networking.Relay(test, db)
			} else if *fullChainFlag {
				fmt.Println("attempting to relay")
				networking.RelayChain(testDesChain, db)
			} else if *hostFlag {
				fmt.Println("attempting to host")
				networking.HostChain(testDesChain, db, *loopFlag)
			}
		}

		if *listenFlag {
			fmt.Println("listening")
			LatestTransaction := networking.ListenRelay()
			fmt.Println(LatestTransaction)

			// Dump fetched tx

			b, err := json.MarshalIndent(LatestTransaction, "", "  ")
			if err != nil {
				fmt.Println("error:", err)
			}
			os.Stdout.Write(b)
		} else if *fetchFlag {
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
	} else if *newChainFlag {
		fmt.Println("creating new chain")

		tsfRef := discovery.NodeID{}

		eDb, err := discovery.NewNodeDatabase(tsfRef, "")

		if err != nil {
			panic(err)
		}

		eDb.WriteDbToMemory(common.GetCurrentDir())

		testcontract := new(contracts.Contract)
		testchain := types.Chain{ParentContract: testcontract, NodeDb: eDb, Version: 0}

		testchain.WriteChainToMemory(common.GetCurrentDir())
	} else if *registerNode {
		common.ThrowWarning("registering node")

		gd, err := networking.GetGateway()

		if err != nil {
			panic(err)
		}

		ip, err := gd.ExternalIP()

		if err != nil {
			panic(err)
		}

		db, err := discovery.ReadDbFromMemory(common.GetCurrentDir())

		if err != nil {
			panic(err)
		}

		ch, err := networking.FetchChain(db)

		if err != nil {
			panic(err)
		}

		db.AddNode(ip, discovery.NodeID{})
		*ch.NodeDb = *db

		networking.RelayChain(ch, db)
	} else {
		common.ThrowWarning("warning: no arguments found")
		fmt.Println("available flags: ")
		flag.PrintDefaults()

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("select flag: ")
		text, _ := reader.ReadString('\n')

		if strings.Contains(text, "relay") {
			*relayFlag = true
		} else if strings.Contains(text, "listen") {
			*listenFlag = true
		} else if strings.Contains(text, "host") {
			*hostFlag = true
		} else if strings.Contains(text, "fetch") {
			*fetchFlag = true
		} else if strings.Contains(text, "new") {
			*newChainFlag = true
		} else if strings.Contains(text, "forever") {
			*loopFlag = true
		} else if strings.Contains(text, "relaychain") {
			*fullChainFlag = true
		} else if strings.Contains(text, "regnode") {
			*registerNode = true
		} else if strings.Contains(text, "noupnp") {
			*noUpNPFlag = true
		}

		if *relayFlag || *listenFlag || *hostFlag || *fetchFlag || *newChainFlag || *loopFlag || *fullChainFlag || *noUpNPFlag || *registerNode {
			main()
		}
	}
}
