package networking

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"time"

	"github.com/mitsukomegumi/indo-go/src/common"
	"github.com/mitsukomegumi/indo-go/src/consensus"
	"github.com/mitsukomegumi/indo-go/src/core/types"
	"github.com/mitsukomegumi/indo-go/src/networking/discovery"

	upnp "github.com/NebulousLabs/go-upnp"
)

const (
	timeout = 15 * time.Second
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
		common.ThrowWarning("verifying tx on current chain")
		fChain := FetchChain(Db)
		if fChain.Transactions[len(fChain.Transactions)-1].InitialWitness.WitnessTime.Before(Tx.InitialWitness.WitnessTime) {
			common.ThrowSuccess("tx passed checks; relaying")
			txBytes := new(bytes.Buffer)
			json.NewEncoder(txBytes).Encode(Tx)
			newConnection(Db.SelfAddr, Db.FindNode(), "relay", txBytes.Bytes()).attempt()
		} else {
			common.ThrowWarning("transaction behind latest chain; fetch latest chain")
		}
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
	if reflect.ValueOf(Ch.NodeDb).IsNil() {
		*Ch = types.Chain{ParentContract: Ch.ParentContract, Identifier: Ch.Identifier, NodeDb: Db, Transactions: Ch.Transactions}
	}
	common.ThrowWarning("attempting to host chain with address " + Ch.NodeDb.SelfAddr)
	if Loop == true {
		for {
			chBytes := new(bytes.Buffer)
			json.NewEncoder(chBytes).Encode(Ch)
			newConnection(Db.SelfAddr, "", "statichostfullchain", chBytes.Bytes()).start(Ch)
		}
	} else {
		chBytes := new(bytes.Buffer)
		json.NewEncoder(chBytes).Encode(Ch)
		newConnection(Db.SelfAddr, "", "statichostfullchain", chBytes.Bytes()).start(Ch)
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

	if err != nil {
		panic(err)
	}

	tempCon.ResolveData(messsage)

	if tempCon.Type == "relay" {
		conn.Close()
		ln.Close()
		return types.DecodeTxFromBytes(tempCon.Data)
	}

	common.ThrowWarning("chain relay found; wanted transaction")
	ln.Close()
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
		conn.Close()
		ln.Close()
		return types.DecodeChainFromBytes(tempCon.Data)
	}

	common.ThrowWarning("transaction relay found; wanted chain")
	conn.Close()
	ln.Close()

	return nil
}

// FetchChain - get current chain from best node; get from nodes with statichostfullchain connection type
func FetchChain(Db *discovery.NodeDatabase) *types.Chain {
	Node := Db.FindNode()

	tempCon := Connection{InitNodeAddr: Db.SelfAddr, DestNodeAddr: Node, Type: "fetchchain"}

	fmt.Println("connection " + tempCon.Type)

	tempCon.AddEvent("started")
	connBytes := new(bytes.Buffer)
	json.NewEncoder(connBytes).Encode(tempCon)

	connec, err := net.Dial("tcp", Node+":3000") // Connect to peer addr

	common.ThrowWarning("attempting to connect to node " + Node + ":3000")

	if err != nil {
		defer func() {
			fmt.Println(err)
		}()
		connec.Close()

		return nil
	}

	connec.Write(connBytes.Bytes())
	fmt.Printf("\n wrote connection meta: %s", connBytes.String())

	connec.SetReadDeadline(time.Now().Add(timeout))

	message, _, err := bufio.NewReader(connec).ReadLine()

	if err != nil {
		common.ThrowWarning("conn err: " + err.Error())
	}

	tempCon.ResolveData(message)

	if tempCon.Type == "statichostfullchain" {
		connec.Close()
		return types.DecodeChainFromBytes(tempCon.Data)
	}

	common.ThrowWarning("chain not found")
	connec.Close()
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
	fChain := FetchChain(Db)
	*Ch = *fChain
	*Ch.NodeDb = *fChain.NodeDb
	Ch.WriteChainToMemory(common.GetCurrentDir())
	Ch.NodeDb.WriteDbToMemory(common.GetCurrentDir())
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

	connec.Close()
}

func (conn *Connection) start(Ch *types.Chain) {
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

	message, _, rErr := bufio.NewReader(connec).ReadLine()

	if rErr != nil {
		common.ThrowWarning(rErr.Error())
	} else {
		tempCon := Connection{}
		tempCon.ResolveData(message)

		fmt.Println("\nConnection type: " + tempCon.Type)

		if tempCon.Type == "fullchain" {
			chain := types.DecodeChainFromBytes(tempCon.Data)
			*Ch = *chain

			common.ThrowSuccess("found chain: ")

			b, err := json.MarshalIndent(chain, "", "  ")
			if err != nil {
				fmt.Println("error:", err)
			}
			os.Stdout.Write(b)

			Ch.WriteChainToMemory(common.GetCurrentDir())
		} else if tempCon.Type == "relay" {
			tx := types.DecodeTxFromBytes(tempCon.Data)
			Ch.AddTransaction(tx)

			common.ThrowSuccess("found transaction: ")

			b, err := json.MarshalIndent(tx, "", "  ")
			if err != nil {
				fmt.Println("error:", err)
			}
			os.Stdout.Write(b)

			Ch.WriteChainToMemory(common.GetCurrentDir())
		} else if tempCon.Type == "fetchchain" {
			fmt.Println("writing to connection")

			fmt.Println(connBytes.Bytes())

			_, wErr := connec.Write(connBytes.Bytes()) // Write connection meta

			if wErr != nil {
				common.ThrowWarning(wErr.Error())
			}
		}
	}

	connec.Close()
	ln.Close()
}

func (conn *Connection) timeout() {
	conn.AddEvent("timed out")
}

func newConnection(initAddr string, destAddr string, connType ConnectionType, data []byte) *Connection {
	if destAddr == "" {
		common.ThrowWarning("\ninitializing new peer connection")
	} else {
		fmt.Printf("forming connection with node %s; ", destAddr)
	}
	fmt.Printf("connection init at %s", common.GetCurrentTime())
	if common.StringInSlice(string(connType), ConnectionTypes) {
		conn := Connection{InitNodeAddr: initAddr, DestNodeAddr: destAddr, Type: connType, Data: data}
		return &conn
	}
	common.ThrowWarning("connection type not valid")
	return nil
}
