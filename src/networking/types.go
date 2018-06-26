package networking

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/mitsukomegumi/indo-go/src/core/types"
	upnp "github.com/nebulouslabs/go-upnp"
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

	Time     time.Time   `json:"inittime"`
	TimeHash *types.Hash `json:"inithash"`

	Type   ConnectionType    `json:"connectiontype"`
	Events []ConnectionEvent `json:"events"`

	Hash *types.Hash `json:"connectionhash"`
}

// ConnectionEvent - string inidicating if event occurred between peers or on network
type ConnectionEvent string

// ConnectionType - represents type of connection being made
type ConnectionType string

// AddEvent - add specified connection event to connection
func (conn *Connection) AddEvent(Event ConnectionEvent) {
	conn.Events = append(conn.Events, Event)
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

// ResolveDataWithChannel - attempts to restore bytes passed via connection to object specified via connectionType
func (conn *Connection) ResolveDataWithChannel(channel chan *[]byte) error {
	val := make([]byte, 100)

	select {
	case tVal := <-channel:
		val = *tVal
	default:
		return errors.New("nil channel")
	}

	plConn := Connection{}
	err := json.NewDecoder(bytes.NewReader(val)).Decode(&plConn)

	if err != nil {
		return err
	}

	plConn.AddEvent("accepted")

	*conn = plConn

	return nil
}

// GetGateway - get reference to current network gateway device
func GetGateway() (*upnp.IGD, error) {
	// connect to router
	d, err := upnp.Discover()
	if err != nil {
		log.Fatal(err)
	}

	return d, err
}
