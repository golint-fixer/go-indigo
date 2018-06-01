package discovery

import (
	"errors"
	"fmt"
	"indogo/src/networking"
	"indogo/src/networking/fastping"
	"net"
	"strings"
	"time"
)

// NodeDatabase - struct holding arrays of IP addresses, node IDs, etc...
type NodeDatabase struct {
	NodeRefDB      []networking.NodeID
	NodePingTimeDB []time.Time
	NodeAddress    []string
	SelfRef        networking.NodeID
}

// NewNodeDatabase - return new node database initialized with self ID
func NewNodeDatabase(selfRef networking.NodeID) *NodeDatabase {
	return &NodeDatabase{SelfRef: selfRef}
}

// AddNode - add specified IP address & ID to node directory
func (db *NodeDatabase) AddNode(ip string, id networking.NodeID) {
	if TestIP(ip) {
		db.NodeAddress = append(db.NodeAddress, ip)
		db.NodePingTimeDB = append(db.NodePingTimeDB, time.Now().UTC())
		db.NodeRefDB = append(db.NodeRefDB, id)
	}
}

// TestIP - ping specified IP address to test for validity
func TestIP(ip string) bool {
	p := fastping.NewPinger()
	ipAddress, err := net.ResolveIPAddr("ip", ip)
	p.AddIPAddr(ipAddress)

	returnVal := false

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
		fmt.Printf("IP %s tested successfully \n", addr.String())
		returnVal = true
	}
	p.OnIdle = func() {
		fmt.Printf("Timed out with IP %s \n", ipAddress)
		returnVal = false
	}
	err = p.Run()
	if err != nil {
		if strings.Contains(err, errors.New("operation not permitted")) {
			fmt.Println("operation requires root priveleges")
		} else {
			fmt.Println(err)
		}
		returnVal = false
	}
	return returnVal
}

// LastPing - Get last ping time for node
func (db *NodeDatabase) LastPing(id networking.NodeID) time.Time {
	nodeIndex := db.GetNodeIndex(id)
	return db.NodePingTimeDB[nodeIndex]
}

// GetNodeIndex - fetch/retrieve node index from node reference
func (db *NodeDatabase) GetNodeIndex(id networking.NodeID) int {
	for k, v := range db.NodeRefDB {
		if id == v {
			return k
		}
	}
	return 0
}
