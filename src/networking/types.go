package networking

// ConnectionTypes - string array representing types of connections that can be
// made on the network, as well as how to resolve them
var ConnectionTypes = []string{"relay", "hostFull"}

// ConnectionEventTypes - preset specifications of acceptable connection event types
var ConnectionEventTypes = []string{"closed", "accepted", "attempted", "started"}

// Connection - struct representing connection between two nodes
type Connection struct {
	InitNodeAddr string `json:"first"`
	DestNodeAddr string `json:"other"`

	Active bool

	connectionData []byte

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
