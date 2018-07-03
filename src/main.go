package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitsukomegumi/indo-go/src/common"
	"github.com/mitsukomegumi/indo-go/src/contracts"
	"github.com/mitsukomegumi/indo-go/src/core/types"
	"github.com/mitsukomegumi/indo-go/src/networking"
	"github.com/mitsukomegumi/indo-go/src/networking/discovery"
	upnp "github.com/nebulouslabs/go-upnp"
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

	var wallet types.Wallet

	if *relayFlag || *listenFlag || *hostFlag || *fetchFlag || *loopFlag || *fullChainFlag || *noUpNPFlag {
		if *listenFlag || *hostFlag {
			common.ThrowWarning("starting host")

			var gd *upnp.IGD

			var err error

			if !*noUpNPFlag {
				common.ThrowWarning("attempting to connect to gateway device")
				gd, err = networking.GetGateway()
			}

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

				var ip string

				var err error

				if gd != nil {
					ip, err = gd.ExternalIP()
					if err != nil {
						panic(err)
					}
				} else {
					ip, err = networking.GetExtIPAddrNoUpNP()
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
			testchain := types.ReadChainFromMemory(common.GetCurrentDir())

			(*testchain).Base = 10

			//Creating new wallet:

			wallet = *types.NewWallet(testchain)

			//Creating witness data:

			signature := types.HexToSignature("4920616d204d697473756b6f204d6567756d69")
			witness := types.NewWitness(1000, signature, 100)

			//Creating transaction, contract, chain

			test := types.NewTransaction(testchain, uint64(1), *wallet.Account, wallet.PrivateKey, wallet.PrivateKeySeeds, types.HexToAddress("4920616d204d697473756b6f204d6567756d69"), common.IntToPointer(0), []byte{0x11, 0x11, 0x11}, nil, nil)

			//Adding witness, transaction to chain

			types.WitnessTransaction(testchain, &wallet, test, &witness)
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
				networking.HostChain(&wallet, testDesChain, &witness, db, *loopFlag)
			}
		}

		if *listenFlag {
			chain := types.ReadChainFromMemory(common.GetCurrentDir())
			wallet := types.NewWallet(chain)

			signature := types.HexToSignature("4920616d204d697473756b6f204d6567756d69")
			witness := types.NewWitness(1000, signature, 100)

			fmt.Println("listening")
			LatestTransaction := networking.ListenRelay(chain, wallet, &witness)
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

/*
	TODO
	- create transaction for reward
	- node registration
	- wallets
*/
