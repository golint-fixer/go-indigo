package networking

import (
	"bufio"
	"bytes"
	"crypto"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/mitsukomegumi/indo-go/src/common"
	"github.com/mitsukomegumi/indo-go/src/consensus"
	"github.com/mitsukomegumi/indo-go/src/core/types"
	"github.com/mitsukomegumi/indo-go/src/networking/discovery"

	upnp "github.com/nebulouslabs/go-upnp"
)

const (
	timeout = 5 * time.Second
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

// GetExtIPAddrNoUpNP - retrieve the external IP address of the current machine w/o upnp
func GetExtIPAddrNoUpNP() (string, error) {
	ip := make([]byte, 100)
	resp, err := http.Get("http://checkip.amazonaws.com/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	_, err = resp.Body.Read(ip)

	if err != nil {
		return "", err
	}

	return string(ip[:len(ip)]), nil
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
func HostChain(Ch *types.Chain, wit *types.Witness, Db *discovery.NodeDatabase, Loop bool) {
	if reflect.ValueOf(Ch.NodeDb).IsNil() {
		*Ch = types.Chain{ParentContract: Ch.ParentContract, Identifier: Ch.Identifier, NodeDb: Db, Transactions: Ch.Transactions, Version: Ch.Version}
	}
	common.ThrowWarning("attempting to host chain with address " + Ch.NodeDb.SelfAddr)
	if Loop == true {
		for {
			chBytes := new(bytes.Buffer)
			json.NewEncoder(chBytes).Encode(Ch)
			newConnection(Db.SelfAddr, "", "statichostfullchain", chBytes.Bytes()).start(Ch, wit)
		}
	} else {
		chBytes := new(bytes.Buffer)
		json.NewEncoder(chBytes).Encode(Ch)
		newConnection(Db.SelfAddr, "", "statichostfullchain", chBytes.Bytes()).start(Ch, wit)
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

	var message bytes.Buffer

	io.Copy(&message, conn)

	if err != nil {
		common.ThrowWarning("conn err: " + err.Error())
	}

	if err != nil {
		panic(err)
	}

	tempCon.ResolveData(message.Bytes())

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

	var message bytes.Buffer

	io.Copy(&message, conn)

	if err != nil {
		common.ThrowWarning("conn err: " + err.Error())
	}

	tempCon.ResolveData(message.Bytes())

	if tempCon.Type == "fullchain" {
		conn.Close()
		ln.Close()

		decodedChain, err := types.DecodeChainFromBytes(tempCon.Data)

		if err != nil {
			panic(err)
		}

		return decodedChain
	}

	common.ThrowWarning("transaction relay found; wanted chain")
	conn.Close()
	ln.Close()

	return nil
}

// FetchChain - get current chain from best node; get from nodes with statichostfullchain connection type
func FetchChain(Db *discovery.NodeDatabase) (*types.Chain, error) {
	Node := Db.FindNode()

	hash := crypto.SHA256.New()

	tempCon := Connection{InitNodeAddr: Db.SelfAddr, DestNodeAddr: Node, Type: "fetchchain", Time: time.Now().UTC(), TimeHash: &types.Hash{}, Hash: &types.Hash{}}

	timeByteArray := hash.Sum([]byte(fmt.Sprintf("%v", tempCon.Time)))

	*tempCon.TimeHash = types.BytesToHash(timeByteArray)

	bArray := hash.Sum([]byte(fmt.Sprintf("%v", tempCon)))

	*tempCon.Hash = types.BytesToHash(bArray)

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
	fmt.Printf("\nwrote connection meta: %s", connBytes.String())

	/*
		finished := make(chan bool)

		go waitForClose(connec, finished)

		<-finished
	*/

	message, err := ioutil.ReadAll(connec)

	if err != nil {
		return nil, err
	}

	if err != nil {
		common.ThrowWarning("conn err: " + err.Error())
		return nil, err
	}

	decomp, err := common.DecompressBytes(message)

	tempCon.ResolveData(decomp)

	if tempCon.Type == "statichostfullchain" {
		connec.Close()

		rCh, err := types.DecodeChainFromBytes(tempCon.Data)

		if err != nil {
			panic(err)
		}

		*Db = *rCh.NodeDb

		decodedChain, err := types.DecodeChainFromBytes(tempCon.Data)

		if err != nil {
			return nil, err
		}

		return decodedChain, nil
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

	//finished := make(chan bool)

	fmt.Println("started")

	/*
		go waitForClose(connec, finished)

		<-finished
	*/

	fmt.Println("finished")

	connec.Close()
	return nil
}

func (conn *Connection) start(Ch *types.Chain, wit *types.Witness) {
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

	data := make(chan []byte, 1000000)

	go handleRequest(connec, data, conn, Ch, wit, finished)

	<-finished

	fmt.Println("finished resolving connection; terminating")

	connec.Close()
	ln.Close()
}

func (conn *Connection) timeout() {
	conn.AddEvent("timed out")
}

func waitForClose(conn net.Conn, finished chan bool) {
	ch := make(chan []byte)
	eCh := make(chan error)

	// Start a goroutine to read from our net connection
	go func(ch chan []byte, eCh chan error) {
		for {
			// try to read the data
			data := make([]byte, 512)
			_, err := conn.Read(data)
			if err != nil {
				// send an error if it's encountered
				eCh <- err
				return
			}
			// send data if we read some.
			ch <- data
		}
	}(ch, eCh)

	ticker := time.Tick(time.Second)
	// continuously read from the connection
	for {
		select {
		// This case means we received data on the connection
		case data := <-ch:
			fmt.Println("Found data: ")

			decomp, err := common.DecompressBytes(data)

			if decomp != nil {
				fmt.Println(string(decomp[:]))

				decodedChain, err := types.DecodeChainFromBytes(decomp)

				fmt.Println(err)

				if decodedChain != nil {
					finished <- true
				}
			} else {
				fmt.Println(err)
			}
		// This case means we got an error and the goroutine has finished
		case err := <-eCh:
			fmt.Println(err)
			break
		// This will timeout on the read.
		case <-ticker:
			common.ThrowWarning("timeout")
			finished <- true
		}
	}
}

func handleRequest(connec net.Conn, data chan []byte, conn *Connection, Ch *types.Chain, wit *types.Witness, finished chan bool) {
	conn.AddEvent("started")
	connBytes := new(bytes.Buffer)
	json.NewEncoder(connBytes).Encode(conn)

	finishedAgainBool := make(chan bool)

	fmt.Println("resolving connection")

	err := resolveConnection(connec, data)
	if err != nil {
		panic(err)
	}

	go finalizeResolvedConnection(data, finishedAgainBool, Ch, wit, connBytes, connec) // call from resolveconnection routine

	<-finishedAgainBool

	finished <- true
}

func resolveConnection(conn net.Conn, buf chan []byte) error {
	err := resolveConnectionData(conn, buf)

	if err != nil {
		return err
	}

	return nil
}

func finalizeResolvedConnection(data chan []byte, finished chan bool, Ch *types.Chain, wit *types.Witness, connBytes *bytes.Buffer, connec net.Conn) {
	tempCon := Connection{}

	rErr := tempCon.ResolveData(<-data)

	if rErr != nil {
		common.ThrowWarning("error while resolving connection data: " + rErr.Error())

		finished <- true
	} else {
		fmt.Println("\nConnection type: " + tempCon.Type)

		if tempCon.Type == "fullchain" {
			chain, err := types.DecodeChainFromBytes(tempCon.Data)

			if err != nil {
				panic(err)
			}

			*Ch = *chain

			common.ThrowSuccess("found chain: ")

			b, err := json.MarshalIndent(chain, "", "  ")
			if err != nil {
				fmt.Println("error:", err)
			}
			os.Stdout.Write(b)

			if !reflect.ValueOf(Ch.NodeDb.NodeAddress).IsNil() {
				common.ThrowSuccess("found node: " + Ch.NodeDb.NodeAddress[len(Ch.NodeDb.NodeAddress)-1])
			}

			Ch.WriteChainToMemory(common.GetCurrentDir())
			Ch.NodeDb.WriteDbToMemory(common.GetCurrentDir())

			finished <- true
		} else if tempCon.Type == "relay" {
			tx := types.DecodeTxFromBytes(tempCon.Data)

			consensus.WitnessTransaction(tx, wit)
			Ch.AddTransaction(tx)

			common.ThrowSuccess("found transaction: ")

			b, err := json.MarshalIndent(tx, "", "  ")
			if err != nil {
				fmt.Println("error:", err)
			}
			os.Stdout.Write(b)

			Ch.WriteChainToMemory(common.GetCurrentDir())

			finished <- true
		} else if tempCon.Type == "fetchchain" {
			fmt.Println("writing to connection")

			b := common.CompressBytes(connBytes.Bytes())

			_, wErr := connec.Write(b) // Write connection meta

			if wErr != nil {
				common.ThrowWarning(wErr.Error())
			}

			finished <- true
		}
		finished <- true
	}
	finished <- true
}

func resolveConnectionData(conn net.Conn, buf chan []byte) error {
	firstLine, _, err := bufio.NewReader(conn).ReadLine()

	if err != nil {
		return err
	}

	err = checkResolution(firstLine)

	if err != nil {
		fmt.Println(err)
		if !strings.Contains(err.Error(), "invalid") && !strings.Contains(err.Error(), "EOF") {
			return err
		}
		var tmpBuffer bytes.Buffer
		io.Copy(&tmpBuffer, conn)

		concatBuf := append(firstLine, tmpBuffer.Bytes()...)

		fmt.Println(string(concatBuf))

		buf <- concatBuf

		return nil
	}

	buf <- firstLine

	return nil
}

func checkResolution(buf []byte) error {
	tempCon := Connection{}
	rErr := tempCon.ResolveData(buf)

	if rErr != nil {
		return rErr
	}
	return nil
}

func newConnection(initAddr string, destAddr string, connType ConnectionType, data []byte) *Connection {
	if destAddr == "" {
		common.ThrowWarning("\ninitializing new peer connection")
	} else {
		fmt.Printf("forming connection with node %s; ", destAddr)
	}
	fmt.Printf("connection init at %s\n", common.GetCurrentTime())
	if common.StringInSlice(string(connType), ConnectionTypes) {
		hash := crypto.SHA256.New()
		conn := Connection{InitNodeAddr: initAddr, DestNodeAddr: destAddr, Type: connType, Data: data, Time: time.Now().UTC(), TimeHash: &types.Hash{}, Hash: &types.Hash{}}

		timeByteArray := hash.Sum([]byte(fmt.Sprintf("%v", conn.Time)))

		*conn.TimeHash = types.BytesToHash(timeByteArray)

		bArray := hash.Sum([]byte(fmt.Sprintf("%v", conn)))

		*conn.Hash = types.BytesToHash(bArray)

		return &conn
	}
	common.ThrowWarning("connection type not valid")
	return nil
}

func isEmpty(channel chan []byte) bool {
	select {
	case msg := <-channel:
		channel <- msg
		return false
	default:
		return true
	}
}
