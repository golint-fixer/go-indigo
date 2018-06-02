package discovery

import (
	"fmt"
	"indo-go/src/common"
	"indo-go/src/networking"
	"indo-go/src/networking/fastping"
	"net"
	"strings"
	"time"
)

const (
	bootStrapNode1 = "108.6.212.149"
)

// NodeDatabase - struct holding arrays of IP addresses, node IDs, etc...
type NodeDatabase struct {
	NodeRefDB      []networking.NodeID
	NodePingTimeDB []time.Time
	NodeAddress    []string
	SelfRef        networking.NodeID
}

// FindNode - find best node to connect to, returns ip address as string
func (db *NodeDatabase) FindNode() string {
	if len(db.NodeAddress) == 0 {
		ReadDbFromMemory(common.GetCurrentDir())
	}
	return db.getBestNode()
}

func (db *NodeDatabase) getBestNode() string {
	x := 0
	bestMatchPingTime := db.NodePingTimeDB[0]
	nodeIndex := 0
	for x != len(db.NodeAddress) {
		if db.NodePingTimeDB[x].After(bestMatchPingTime) {
			bestMatchPingTime = db.NodePingTimeDB[x]
			nodeIndex = x
		}
		x++
	}
	return db.NodeAddress[nodeIndex]
}

// NewNodeDatabase - return new node database initialized with self ID
func NewNodeDatabase(selfRef networking.NodeID) *NodeDatabase {
	readDb := ReadDbFromMemory(common.GetCurrentDir())
	if readDb != nil {
		fmt.Println("read existing node database from mem")
		return readDb
	}
	return &NodeDatabase{SelfRef: selfRef}
}

// AddNode - add specified IP address & ID to node directory
func (db *NodeDatabase) AddNode(ip string, id networking.NodeID) {
	if !strings.Contains(ip, "192.") {
		if TestIP(ip) {
			fmt.Println("adding node to database")
			db.NodeAddress = append(db.NodeAddress, ip)
			db.NodePingTimeDB = append(db.NodePingTimeDB, time.Now().UTC())
			db.NodeRefDB = append(db.NodeRefDB, id)
		}
	} else {
		common.ThrowWarning("database error: node cannot be internal")
	}
}

// WriteDbToMemory - create serialized instance of specified NodeDatabase in specified path (string)
func (db *NodeDatabase) WriteDbToMemory(path string) {
	err := common.WriteGob(path+"nodeDb.gob", db)

	if err != nil {
		fmt.Println(err)
	} else {
		common.ThrowSuccess("\nobject written to memory")
	}
}

// ReadDbFromMemory - read serialized object of specified node database from specified path
func ReadDbFromMemory(path string) *NodeDatabase {
	tempDb := new(NodeDatabase)

	error := common.ReadGob(path+"nodeDb.gob", tempDb)
	if error != nil {
		fmt.Println(error)
	} else {
		return tempDb
	}
	return nil
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
		if returnVal != true {
			fmt.Printf("Timed out with IP %s \n", ipAddress)
			returnVal = false
		}
	}
	err = p.Run()
	if err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
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
