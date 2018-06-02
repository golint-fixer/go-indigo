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
	"indo-go/src/networking/upnp"
	"net"
)

// AddPortMapping - add port mapping on specified port
func AddPortMapping(port int) {
	mapping := new(upnp.Upnp)
	if err := mapping.AddPortMapping(port, port, "TCP"); err == nil {
		fmt.Println("port mapping added")
	} else {
		fmt.Printf("port mapping failed with err %s", err)
	}
}

// Relay - push localized or received transaction to further node.
func Relay(Tx *types.Transaction, Db *discovery.NodeDatabase) {
	AddPortMapping(3000)
	txBytes := new(bytes.Buffer)
	json.NewEncoder(txBytes).Encode(Tx)
	newConnection(Db.SelfAddr, Db.FindNode(), "relay", txBytes.Bytes()).attempt()
}

// ListenRelay - listen for transaction relays, relay to full node or host
func ListenRelay() *types.Transaction {
	AddPortMapping(3000)

	tempCon := Connection{}

	ln, err := net.Listen("tcp", ":3000")

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	conn, err := ln.Accept()

	messsage, _, err := bufio.NewReader(conn).ReadLine()
	tempCon.ResolveData(messsage)
	return types.DecodeTxFromBytes(tempCon.Data)
}

// ListenRelayWithAdd - listen for transaction relays, add to local chain
func ListenRelayWithAdd(Ch *types.Chain, Wit *types.Witness, Db *discovery.NodeDatabase) {
	tx := ListenRelay()
	consensus.WitnessTransaction(tx, Wit)
	Ch.AddTransaction(tx)
	Ch.WriteChainToMemory(common.GetCurrentDir())
	Relay(tx, Db)
}

func handleReceivedBytes(b []byte) *Connection {
	tempConn := Connection{}
	tempConn.ResolveData(b)
	return &tempConn
}

// TODO: encode

func (conn *Connection) attempt() {

	connBytes := new(bytes.Buffer)
	json.NewEncoder(connBytes).Encode(conn)

	connec, err := net.Dial("tcp", conn.DestNodeAddr+":3000") // Connect to peer addr
	connec.Write(connBytes.Bytes())                           // Write connection meta
	if err != nil {
		fmt.Println(err)
	} else {
		conn.AddEvent("started")
	}
}

func newConnection(initAddr string, destAddr string, connType ConnectionType, data []byte) *Connection {
	if common.StringInSlice(string(connType), ConnectionTypes) {
		conn := Connection{InitNodeAddr: initAddr, DestNodeAddr: destAddr, Type: connType, Data: data}
		return &conn
	}
	common.ThrowWarning("connection type not valid")
	return nil
}

//https://systembash.com/a-simple-go-tcp-server-and-tcp-client/
