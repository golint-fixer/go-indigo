package networking

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/mitsukomegumi/indo-go/src/common"
	"github.com/mitsukomegumi/indo-go/src/consensus"
	"github.com/mitsukomegumi/indo-go/src/core/types"
	"github.com/mitsukomegumi/indo-go/src/networking/discovery"

	upnp "github.com/NebulousLabs/go-upnp"
)

const (
	timeout = 15 * time.Second
	rDelay  = 2 * time.Second
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
func Relay(Tx *types.Transaction, Db *discovery.NodeDatabase) error {
	if !reflect.ValueOf(Tx.InitialWitness).IsNil() {
		common.ThrowWarning("verifying tx on current chain")
		fChain, err := FetchChain(Db)

		if err != nil {
			return err
		}

		if fChain.Transactions[len(fChain.Transactions)-1].InitialWitness.WitnessTime.Before(Tx.InitialWitness.WitnessTime) {
			common.ThrowSuccess("tx passed checks; relaying")
			txBytes := new(bytes.Buffer)
			json.NewEncoder(txBytes).Encode(Tx)
			newConnection(Db.SelfAddr, Db.FindNode(), "relay", txBytes.Bytes()).attempt()

			return nil
		}
		return errors.New("transaction behind latest chain; fetch latest chain")
	}
	return errors.New("operation not permitted; transaction not witnessed")
}

// RelayChain - push localized or received chain to further node
func RelayChain(Ch *types.Chain, Db *discovery.NodeDatabase) error {
	chBytes := new(bytes.Buffer)
	err := json.NewEncoder(chBytes).Encode(Ch)

	if err != nil {
		return err
	}

	newConnection(Db.SelfAddr, Db.FindNode(), "fullchain", chBytes.Bytes()).attempt()

	return nil
}

// HostChain - host localized chain to forwarded port
func HostChain(Ch *types.Chain, Db *discovery.NodeDatabase, Loop bool) {
	if reflect.ValueOf(Ch.NodeDb).IsNil() {
		*Ch = types.Chain{ParentContract: Ch.ParentContract, Identifier: Ch.Identifier, NodeDb: Db, Transactions: Ch.Transactions, Version: Ch.Version}
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

	message, err := ioutil.ReadAll(conn)

	if err != nil {
		common.ThrowWarning("conn err: " + err.Error())
	}

	if err != nil {
		panic(err)
	}

	tempCon.ResolveData(message)

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

	message, err := ioutil.ReadAll(conn)

	if err != nil {
		common.ThrowWarning("conn err: " + err.Error())
	}

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
func FetchChain(Db *discovery.NodeDatabase) (*types.Chain, error) {
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

		return nil, err
	}

	connec.Write(connBytes.Bytes())
	fmt.Printf("\n wrote connection meta: %s", connBytes.String())

	message, err := ioutil.ReadAll(connec)

	if err != nil {
		common.ThrowWarning("conn err: " + err.Error())
		return nil, err
	}

	tempCon.ResolveData(common.DecompressBytes(message))

	if tempCon.Type == "statichostfullchain" {
		connec.Close()

		rCh := types.DecodeChainFromBytes(tempCon.Data)

		*Db = *rCh.NodeDb

		return types.DecodeChainFromBytes(tempCon.Data), nil
	}

	common.ThrowWarning("chain not found")
	connec.Close()
	return nil, errors.New("chain not found")
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
func FetchChainWithAdd(Ch *types.Chain, Db *discovery.NodeDatabase) error {
	fChain, err := FetchChain(Db)

	if err != nil {
		return err
	}

	*Ch = *fChain
	*Ch.NodeDb = *fChain.NodeDb
	Ch.WriteChainToMemory(common.GetCurrentDir())
	Ch.NodeDb.WriteDbToMemory(common.GetCurrentDir())

	return nil
}

func handleReceivedBytes(b []byte) *Connection {
	tempConn := Connection{}
	tempConn.ResolveData(b)
	return &tempConn
}

func (conn *Connection) attempt() error {
	conn.AddEvent("started")
	connBytes := new(bytes.Buffer)
	json.NewEncoder(connBytes).Encode(conn)

	common.ThrowWarning("\nattempting to dial address: " + conn.DestNodeAddr + ":3000")

	connec, err := net.Dial("tcp", conn.DestNodeAddr+":3000") // Connect to peer addr

	if err != nil {
		return err
	}

	connec.SetDeadline(time.Now().Add(timeout)) // Set timeout
	connec.Write(connBytes.Bytes())             // Write connection meta

	fmt.Println(connBytes.Bytes())

	conn.AddEvent("started")

	connec.Close()

	return nil
}

func (conn *Connection) start(Ch *types.Chain) {
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

	finished := make(chan bool)

	go handleRequest(connec, conn, Ch, finished)

	<-finished

	//connec.Close()
	ln.Close()
}

func (conn *Connection) timeout() {
	conn.AddEvent("timed out")
}

func handleRequest(connec net.Conn, conn *Connection, Ch *types.Chain, finished chan bool) {
	conn.AddEvent("started")
	connBytes := new(bytes.Buffer)
	json.NewEncoder(connBytes).Encode(conn)

	buf := new([]byte)

	*buf = resolveConnection(connec)

	tempCon := Connection{}
	rErr := tempCon.ResolveData(*buf)

	fmt.Println(buf)

	if rErr != nil {
		common.ThrowWarning(rErr.Error())
	} else {
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

			common.ThrowSuccess("found node: " + Ch.NodeDb.NodeAddress[len(Ch.NodeDb.NodeAddress)-1])

			Ch.WriteChainToMemory(common.GetCurrentDir())
			Ch.NodeDb.WriteDbToMemory(common.GetCurrentDir())
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

			b := common.CompressBytes(connBytes.Bytes())

			_, wErr := connec.Write(b) // Write connection meta

			if wErr != nil {
				common.ThrowWarning(wErr.Error())
			}
		}
	}
}

func resolveConnection(conn net.Conn) []byte {
	testResolution, err := resolveSimple(conn)

	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			overflowResolution, err := resolveOverflow(conn)

			if err != nil {
				panic(err)
			}
			return overflowResolution
		}
		panic(err)
	}

	return testResolution
}

func resolveOverflow(conn net.Conn) ([]byte, error) {
	data, err := ioutil.ReadAll(conn)

	if err != nil {
		panic(err)
	}

	return data, nil
}

func resolveSimple(conn net.Conn) ([]byte, error) {
	data, _, err := bufio.NewReader(conn).ReadLine()

	if err != nil {
		return nil, err
	}

	return data, nil
}

func newConnection(initAddr string, destAddr string, connType ConnectionType, data []byte) *Connection {
	if destAddr == "" {
		common.ThrowWarning("\ninitializing new peer connection")
	} else {
		fmt.Printf("forming connection with node %s; ", destAddr)
	}
	fmt.Printf("connection init at %s\n", common.GetCurrentTime())
	if common.StringInSlice(string(connType), ConnectionTypes) {
		conn := Connection{InitNodeAddr: initAddr, DestNodeAddr: destAddr, Type: connType, Data: data}
		return &conn
	}
	common.ThrowWarning("connection type not valid")
	return nil
}
