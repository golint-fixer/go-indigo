package networking

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"
	"time"

	"github.com/mitsukomegumi/indo-go/src/common"
	"github.com/mitsukomegumi/indo-go/src/consensus"
	"github.com/mitsukomegumi/indo-go/src/core/types"
	"github.com/mitsukomegumi/indo-go/src/networking/discovery"

	upnp "github.com/NebulousLabs/go-upnp"
)

const (
	timeout = 5 * time.Second
)

func forward(GatewayDevice *upnp.IGD) {
	// discover external IP
	ip, err := GatewayDevice.ExternalIP()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("current node external ip:", ip)

	// forward a port
	err = GatewayDevice.Forward(3000, "resourceforwarding")
	if err != nil {
		log.Fatal(err)
	}
}

func removeMapping(GatewayDevice *upnp.IGD) {
	// discover external IP
	ip, err := GatewayDevice.ExternalIP()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("current node external ip:", ip)

	// remove port mappings
	err = GatewayDevice.Clear(3000)
	if err != nil {
		log.Fatal(err)
	}
}

// PrepareForConnection - forward all necessary ports to decrease redundant speed limitations
func PrepareForConnection(GatewayDevice *upnp.IGD, db *discovery.NodeDatabase) {
	forward(GatewayDevice)
}

// DisableConnections - remove all necessary port mappings
func DisableConnections(GatewayDevice *upnp.IGD, db *discovery.NodeDatabase) {
	removeMapping(GatewayDevice)
}

// GetExtIPAddr - retrieve the external IP address of the current machine
func GetExtIPAddr() string {
	// connect to router
	d, err := upnp.Discover()
	if err != nil {
		log.Fatal(err)
	}

	// discover external IP
	ip, err := d.ExternalIP()
	if err != nil {
		log.Fatal(err)
	}
	return ip
}

// Relay - push localized or received transaction to further node
func Relay(Tx *types.Transaction, Db *discovery.NodeDatabase) {
	if !reflect.ValueOf(Tx.InitialWitness).IsNil() {
		//if ListenChain().Transactions[len(ListenChain().Transactions)-1].InitialWitness.WitnessTime.Before(Tx.InitialWitness.WitnessTime) { // Causes infinite loop if no nodes serving chain
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
	chBytes := new(bytes.Buffer)
	json.NewEncoder(chBytes).Encode(Ch)
	newConnection(Db.SelfAddr, Db.FindNode(), "fullchain", chBytes.Bytes()).attempt()
}

// HostChain - host localized chain to forwarded port
func HostChain(Ch *types.Chain, Db *discovery.NodeDatabase, Loop bool) {
	if Loop == true {
		for {
			chBytes := new(bytes.Buffer)
			json.NewEncoder(chBytes).Encode(Ch)
			newConnection(Db.SelfAddr, "", "statichostfullchain", chBytes.Bytes()).start()
		}
	} else {
		chBytes := new(bytes.Buffer)
		json.NewEncoder(chBytes).Encode(Ch)
		newConnection(Db.SelfAddr, "", "statichostfullchain", chBytes.Bytes()).start()
	}
}

// ListenRelay - listen for transaction relays, relay to full node or host
func ListenRelay() *types.Transaction {
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

	messsage, _, err := bufio.NewReader(conn).ReadLine()
	tempCon.ResolveData(messsage)

	if tempCon.Type == "relay" {
		return types.DecodeTxFromBytes(tempCon.Data)
	}

	common.ThrowWarning("chain relay found; wanted transaction")
	conn.Close()

	return nil
}

// ListenChain - listen for chain relays, relay to full node or host
func ListenChain() *types.Chain {
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
	conn.Close()

	return nil
}

// FetchChain - get current chain from best node; get from nodes with statichostfullchain connection type
func FetchChain(Db *discovery.NodeDatabase) *types.Chain {
	tempCon := Connection{}

	connec, err := net.Dial("tcp", Db.FindNode()+":3000") // Connect to peer addr
	connec.SetDeadline(time.Now().Add(timeout))

	if err != nil {
		defer func() {
			fmt.Println(err)
		}()
		connec.Close()

		return nil
	}
	message, _, err := bufio.NewReader(connec).ReadLine()

	tempCon.ResolveData(message)

	if tempCon.Type == "statichostfullchain" {
		return types.DecodeChainFromBytes(tempCon.Data)
	}

	common.ThrowWarning("chain not found")
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

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	connec.Write(connBytes.Bytes()) // Write connection meta
	connec.Close()
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
