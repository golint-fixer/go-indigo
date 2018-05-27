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
func NewNodeDatabase(selfRef networking.NodeID) *nodeDatabase {
	return &nodeDatabase{SelfRef: selfRef}
}

func (db *nodeDatabase) addNode(ip string, id networking.NodeID) {
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr(ip, os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Node tested successfully: " + ip)
	}
}

func (db *nodeDatabase) lastPing(id networking.NodeID) time.Time {
	nodeIndex := db.getNodeIndex(id)
	return db.NodePingTimeDB[nodeIndex]
}

func (db *nodeDatabase) getNodeIndex(id networking.NodeID) int {
	for k, v := range db.NodeRefDB {
		if id == v {
			return k
		}
	}
	return 0
}
