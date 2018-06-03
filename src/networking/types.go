package networking

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// ConnectionTypes - string array representing types of connections that can be
// made on the network, as well as how to resolve them
var ConnectionTypes = []string{"relay", "fullchain", "statichost", "statichostfullchain"}

// ConnectionEventTypes - preset specifications of acceptable connection event types
var ConnectionEventTypes = []string{"closed", "accepted", "attempted", "started"}

// Connection - struct representing connection between two nodes
type Connection struct {
	InitNodeAddr string `json:"first"`
	DestNodeAddr string `json:"other"`

	Active bool

	Data []byte

	Type   ConnectionType
	Events []ConnectionEvent
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
func (conn *Connection) ResolveData(b []byte) {
	err := binary.Read(bytes.NewBuffer(b[:]), binary.BigEndian, *conn)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	conn.AddEvent("accepted")
}
