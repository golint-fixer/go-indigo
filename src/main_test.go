package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/mitsukomegumi/indo-go/src/common"
	"github.com/mitsukomegumi/indo-go/src/contracts"
	"github.com/mitsukomegumi/indo-go/src/core/types"
	"github.com/mitsukomegumi/indo-go/src/networking"
	"github.com/mitsukomegumi/indo-go/src/networking/discovery"
)

func TestNewChain(t *testing.T) {
	err := NewChain()

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestRelayTx(t *testing.T) {
	tsfRef := discovery.NodeID{}

	db, err := discovery.NewNodeDatabase(tsfRef, "")

	if err != nil {
		t.Errorf("Node database creation failed: %s", err.Error())
	}

	//Creating transaction, contract, chain

	eDb, err := discovery.NewNodeDatabase(tsfRef, "")

	if err != nil {
		t.Errorf("Node database creation failed: %s", err.Error())
	}

	wErr := eDb.WriteDbToMemory(common.GetCurrentDir())

	if wErr != nil {
		t.Errorf("Node database serialization failed: %s", err.Error())
	}

	testcontract := new(contracts.Contract)
	testchain := types.Chain{ParentContract: testcontract, NodeDb: eDb, Version: 0}
	sErr := testchain.WriteChainToMemory(common.GetCurrentDir())

	//Creating new account:

	wallet := *types.NewWallet(&testchain)

	//Creating witness data:

	signature := types.HexToSignature("4920616d204d697473756b6f204d6567756d69")
	witness := types.NewWitness(1000, signature, 100)

	if sErr != nil {
		t.Errorf("Chain serialization failed: %s", sErr.Error())
	}

	test := types.NewTransaction(&testchain, uint64(1), *wallet.Account, wallet.PrivateKey, wallet.PrivateKeySeeds, types.HexToAddress("4920616d204d697473756b6f204d6567756d69"), common.IntToPointer(1000), []byte{0x11, 0x11, 0x11}, nil, nil)

	//Adding witness, transaction to chain

	types.WitnessTransaction(&testchain, &wallet, test, &witness)
	testchain.AddTransaction(test)

	//Test chain serialization

	fmt.Println("attempting to relay")
	rErr := networking.Relay(test, db)

	if rErr != nil {
		t.Errorf(rErr.Error())
	}
}

func TestRelayChain(t *testing.T) {
	NewChain()
	chain := types.ReadChainFromMemory(common.GetCurrentDir())
	db, err := discovery.ReadDbFromMemory(common.GetCurrentDir())

	if err != nil {
		t.Errorf("Node database deserialization failed: %s", err.Error())
	}

	rErr := networking.RelayChain(chain, db)

	if rErr != nil {
		t.Errorf("Chain relay failed: %s", err.Error())
	}
}

func TestFetchChain(t *testing.T) {
	selfRef := discovery.NodeID{}
	selfAddr := ""

	db, err := discovery.NewNodeDatabase(selfRef, selfAddr)

	if err != nil {
		t.Errorf(err.Error())
	}

	testDesChain := types.Chain{}
	fmt.Println("attempting to fetch chain")
	err = networking.FetchChainWithAdd(&testDesChain, db)

	if err != nil {
		t.Errorf(err.Error())
	}

	// Dump fetched chain

	b, err := json.MarshalIndent(testDesChain, "", "  ")
	if err != nil {
		t.Errorf(err.Error())
	}
	os.Stdout.Write(b)
}

func NewChain() error {
	tsfRef := discovery.NodeID{}

	eDb, err := discovery.NewNodeDatabase(tsfRef, "")

	if err != nil {
		return err
	}

	wErr := eDb.WriteDbToMemory(common.GetCurrentDir())

	if wErr != nil {
		return wErr
	}

	testcontract := new(contracts.Contract)
	testchain := types.Chain{ParentContract: testcontract, NodeDb: eDb, Version: 0}
	sErr := testchain.WriteChainToMemory(common.GetCurrentDir())

	if sErr != nil {
		return sErr
	}

	return nil
}
