package discovery

import (
	"fmt"
	"indogo/src/networking"
	"indogo/src/networking/fastping"
	"net"
	"time"
)

// NodeDatabase - struct holding arrays of IP addresses, node IDs, etc...
type NodeDatabase struct {
	NodeRefDB      []networking.NodeID
	NodePingDB     []int
	NodePingTimeDB []time.Time
	NodeAddress    []string
	SelfRef        networking.NodeID
}

// NewNodeDatabase - return new node database initialized with self ID
func NewNodeDatabase(selfRef networking.NodeID) *NodeDatabase {
	return &NodeDatabase{SelfRef: selfRef}
}

// AddNode - add specified IP address & ID to node directory
func (db *NodeDatabase) AddNode(ip string) {
	p := fastping.NewPinger()
	ipAddress := net.IPAddr{IP: net.IP([]byte(ip))}
	p.AddIPAddr(&ipAddress)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
		db.NodeAddress = append(db.NodeAddress, ip)
	}
	p.OnIdle = func() {
		fmt.Printf("Timed out with IP %s", ipAddress)
	}
	err := p.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func (db *NodeDatabase) lastPing(id networking.NodeID) time.Time {
	nodeIndex := db.getNodeIndex(id)
	return db.NodePingTimeDB[nodeIndex]
}

func (db *NodeDatabase) getNodeIndex(id networking.NodeID) int {
	for k, v := range db.NodeRefDB {
		if id == v {
			return k
		}
	}
	return 0
}
