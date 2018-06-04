package networking

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"indo-go/src/common"
	"indo-go/src/consensus"
	"indo-go/src/core/types"
	"indo-go/src/networking/discovery"
	"net"
	"reflect"
	"time"
)

const (
	timeout = 5
)

// Relay - push localized or received transaction to further node
func Relay(Tx *types.Transaction, Db *discovery.NodeDatabase) {
	if !reflect.ValueOf(Tx.InitialWitness).IsNil() {
		//if ListenChain().Transactions[len(ListenChain().Transactions)-1].InitialWitness.WitnessTime.Before(Tx.InitialWitness.WitnessTime) { // Causes infinite loop if no nodes serving chain
		//AddPortMapping(3000)
		common.ThrowSuccess("tx passed checks; relaying")
		txBytes := new(bytes.Buffer)
		json.NewEncoder(txBytes).Encode(Tx)
		newConnection(Db.SelfAddr, Db.FindNode(), "relay", txBytes.Bytes()).attempt()
		//} else {
		//common.ThrowWarning("transaction behind latest chain; fetch latest chain")
		//}
	} else {
		common.ThrowWarning("operation not permitted; transaction not witness")
	}
}

// RelayChain - push localized or received chain to further node
func RelayChain(Ch *types.Chain, Db *discovery.NodeDatabase) {
	//AddPortMapping(3000)
	chBytes := new(bytes.Buffer)
	json.NewEncoder(chBytes).Encode(Ch)
	newConnection(Db.SelfAddr, Db.FindNode(), "fullchain", chBytes.Bytes()).attempt()
}

// HostChain - host localized chain to forwarded port
func HostChain(Ch *types.Chain, Db *discovery.NodeDatabase) {
	//AddPortMapping(3000)
	chBytes := new(bytes.Buffer)
	json.NewEncoder(chBytes).Encode(Ch)
	newConnection(Db.SelfAddr, "", "statichostfullchain", chBytes.Bytes()).start()
}

// ListenRelay - listen for transaction relays, relay to full node or host
func ListenRelay() *types.Transaction {
	//AddPortMapping(3000)

	tempCon := Connection{}

	ln, err := net.Listen("tcp", ":3000")

	if err != nil {
		fmt.Println(err)
		//panic(err)
	}

	conn, err := ln.Accept()
	conn.SetDeadline(time.Now().Add(timeout))

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	messsage, _, err := bufio.NewReader(conn).ReadLine()
	tempCon.ResolveData(messsage)

	if tempCon.Type == "relay" {
		return types.DecodeTxFromBytes(tempCon.Data)
	}

	common.ThrowWarning("chain relay found; wanted transaction")

	return nil
}

// ListenChain - listen for chain relays, relay to full node or host
func ListenChain() *types.Chain {
	//AddPortMapping(3000)

	tempCon := Connection{}

	ln, err := net.Listen("tcp", ":3000")

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	conn, err := ln.Accept()
	conn.SetDeadline(time.Now().Add(timeout))

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	message, _, err := bufio.NewReader(conn).ReadLine()
	tempCon.ResolveData(message)

	if tempCon.Type == "fullchain" {
		return types.DecodeChainFromBytes(tempCon.Data)
	}

	common.ThrowWarning("transaction relay found; wanted chain")

	return nil
}

// FetchChain - get current chain from best node; get from nodes with statichostfullchain connection type
func FetchChain(Db *discovery.NodeDatabase) *types.Chain {
	connec, err := net.Dial("tcp", Db.FindNode()+":3000") // Connect to peer addr
	connec.SetDeadline(time.Now().Add(timeout))

	tempCon := Connection{}

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	message, _, err := bufio.NewReader(connec).ReadLine()
	tempCon.ResolveData(message)

	if tempCon.Type == "statichostfullchain" {
		return types.DecodeChainFromBytes(tempCon.Data)
	}

	common.ThrowWarning("chain host not found")

	return nil
}

// ListenRelayWithAdd - listen for transaction relays, add to local chain
func ListenRelayWithAdd(Ch *types.Chain, Wit *types.Witness, Db *discovery.NodeDatabase) {
	tx := ListenRelay()
	consensus.WitnessTransaction(tx, Wit)
	Ch.AddTransaction(tx)
	Ch.WriteChainToMemory(common.GetCurrentDir())
	Relay(tx, Db)
}

// ListenChainWithAdd - listen for chain relays, set local chain to result
func ListenChainWithAdd(Ch *types.Chain, Db *discovery.NodeDatabase) {
	*Ch = *ListenChain()
	Ch.WriteChainToMemory(common.GetCurrentDir())
	RelayChain(Ch, Db)
}

// FetchChainWithAdd - fetch chain, set local chain to result
func FetchChainWithAdd(Ch *types.Chain, Db *discovery.NodeDatabase) {
	*Ch = *FetchChain(Db)
	Ch.WriteChainToMemory(common.GetCurrentDir())
}

func handleReceivedBytes(b []byte) *Connection {
	tempConn := Connection{}
	tempConn.ResolveData(b)
	return &tempConn
}

func (conn *Connection) attempt() {
	conn.AddEvent("started")
	connBytes := new(bytes.Buffer)
	json.NewEncoder(connBytes).Encode(conn)

	common.ThrowWarning("attempting to dial address: " + conn.DestNodeAddr + ":3000")

	connec, err := net.Dial("tcp", conn.DestNodeAddr+":3000") // Connect to peer addr
	connec.SetDeadline(time.Now().Add(timeout))               // Set timeout
	connec.Write(connBytes.Bytes())                           // Write connection meta

	if err != nil {
		fmt.Println(err)
	} else {
		conn.AddEvent("started")
	}
}

func (conn *Connection) start() {
	conn.AddEvent("started")
	connBytes := new(bytes.Buffer)
	json.NewEncoder(connBytes).Encode(conn)

	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Println(err) // Print panic meta
		panic(err)       // Panic
	}

	connec, err := ln.Accept() // Accept peer connection
	connec.SetDeadline(time.Now().Add(timeout))

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	connec.Write(connBytes.Bytes()) // Write connection meta
}

func (conn *Connection) timeout() {
	conn.AddEvent("timed out")
}

func newConnection(initAddr string, destAddr string, connType ConnectionType, data []byte) *Connection {
	fmt.Printf("forming connection with node %s; ", destAddr)
	if common.StringInSlice(string(connType), ConnectionTypes) {
		conn := Connection{InitNodeAddr: initAddr, DestNodeAddr: destAddr, Type: connType, Data: data}
		return &conn
	}
	common.ThrowWarning("connection type not valid")
	return nil
}
