package types

import (
	"bytes"
	"crypto"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"reflect"
	"time"

	"github.com/mitsukomegumi/indo-go/src/common"
	"github.com/mitsukomegumi/indo-go/src/networking/discovery"
)

const (
	timeout = 5 * time.Second
	rDelay  = 2 * time.Second
)

// ConnectionTypes - string array representing types of connections that can be
// made on the network, as well as how to resolve them
var ConnectionTypes = []string{"relay", "fullchain", "statichost", "statichostfullchain", "fetchchain"}

// ConnectionEventTypes - preset specifications of acceptable connection event types
var ConnectionEventTypes = []string{"closed", "accepted", "attempted", "started", "timed out"}

// Connection - struct representing connection between two nodes
type Connection struct {
	InitNodeAddr string `json:"first"`
	DestNodeAddr string `json:"other"`

	Active bool `json:"active"`

	Data []byte `json:"data"`

	Time     time.Time `json:"inittime"`
	TimeHash *Hash     `json:"inithash"`

	Type   ConnectionType    `json:"connectiontype"`
	Events []ConnectionEvent `json:"events"`

	Extra []byte `json:"extradata"`

	Hash *Hash `json:"connectionhash"`
}

// ConnectionEvent - string inidicating if event occurred between peers or on network
type ConnectionEvent string

// ConnectionType - represents type of connection being made
type ConnectionType string

// AddEvent - add specified connection event to connection
func (conn *Connection) AddEvent(Event ConnectionEvent) {
	conn.Events = append(conn.Events, Event)
}

// Relay - push localized or received transaction to further node
func Relay(Tx *Transaction, Db *discovery.NodeDatabase) error {
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
			time.Sleep(20 * time.Millisecond)
			newConnection(Db.SelfAddr, Db.FindNode(), "relay", txBytes.Bytes()).attempt()

			return nil
		}
		return errors.New("transaction behind latest chain; fetch latest chain")
	}
	return errors.New("operation not permitted; transaction not witnessed")
}

// FetchChain - get current chain from best node; get from nodes with statichostfullchain connection type
func FetchChain(Db *discovery.NodeDatabase) (*Chain, error) {
	Node := Db.FindNode()

	hash := crypto.SHA256.New()

	tempCon := Connection{InitNodeAddr: Db.SelfAddr, DestNodeAddr: Node, Type: "fetchchain", Time: time.Now().UTC(), TimeHash: &Hash{}, Hash: &Hash{}}

	timeByteArray := hash.Sum([]byte(fmt.Sprintf("%v", tempCon.Time)))

	*tempCon.TimeHash = BytesToHash(timeByteArray)

	bArray := hash.Sum([]byte(fmt.Sprintf("%v", tempCon)))

	*tempCon.Hash = BytesToHash(bArray)

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

		rCh, err := DecodeChainFromBytes(tempCon.Data)

		if err != nil {
			panic(err)
		}

		*Db = *rCh.NodeDb

		decodedChain, err := DecodeChainFromBytes(tempCon.Data)

		if err != nil {
			return nil, err
		}

		return decodedChain, nil
	}

	common.ThrowWarning("chain not found")
	connec.Close()
	return nil, errors.New("chain not found")
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
		conn := Connection{InitNodeAddr: initAddr, DestNodeAddr: destAddr, Type: connType, Data: data, Time: time.Now().UTC(), TimeHash: &Hash{}, Hash: &Hash{}}

		timeByteArray := hash.Sum([]byte(fmt.Sprintf("%v", conn.Time)))

		*conn.TimeHash = BytesToHash(timeByteArray)

		bArray := hash.Sum([]byte(fmt.Sprintf("%v", conn)))

		*conn.Hash = BytesToHash(bArray)

		return &conn
	}
	common.ThrowWarning("connection type not valid")
	return nil
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

// ResolveData - attempts to restore bytes passed via connection to object specified via connectionType
func (conn *Connection) ResolveData(b []byte) error {
	plConn := Connection{}
	err := json.NewDecoder(bytes.NewReader(b)).Decode(&plConn)

	if err != nil {
		return err
	}

	plConn.AddEvent("accepted")

	*conn = plConn

	return nil
}
