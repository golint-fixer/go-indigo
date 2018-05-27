package discovery

import (
	"fmt"
	"indogo/src/networking"
	"indogo/src/networking/fastping"
	"net"
	"os"
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
	ra, err := net.ResolveIPAddr(ip, os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Node tested successfully: " + ip)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
	}
	p.OnIdle = func() {
		fmt.Println("finish")
	}
	err = p.Run()
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
