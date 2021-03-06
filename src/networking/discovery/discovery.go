package discovery

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/mitsukomegumi/indo-go/src/common"
	"github.com/mitsukomegumi/indo-go/src/networking/fastping"
)

const (
	bootStrapNode1Addr = "10.144.4.68"
)

// NodeDatabase - struct holding arrays of IP addresses, node IDs, etc...
type NodeDatabase struct {
	NodeRefDB          []NodeID
	NodePingTimeDB     []time.Time
	NodeAddress        []string
	SelfRef            NodeID
	SelfAddr           string
	BootstrapNodeAddrs []string
}

// NodeID - byte array identifying individual node
type NodeID [64]byte

// FindNode - find best node to connect to, returns ip address as string
func (db *NodeDatabase) FindNode() string {
	if !reflect.ValueOf(db).IsNil() {
		if len(db.NodeAddress) == 0 {
			ReadDbFromMemory(common.GetCurrentDir())
		}
		return db.getBestNode()
	}
	common.ThrowWarning("nil db")
	return "10.144.4.68"
}

func (db *NodeDatabase) getBestNode() string {
	if len(db.NodeAddress) > 0 {
		x := 0
		bestMatchPingTime := db.NodePingTimeDB[0]
		nodeIndex := 0
		for x != len(db.NodeAddress)-1 {
			if db.NodePingTimeDB[x].After(bestMatchPingTime) {
				bestMatchPingTime = db.NodePingTimeDB[x]
				nodeIndex = x
			}
			x++
		}
		return db.NodeAddress[nodeIndex]
	}
	return db.getBootstrap()
}

func (db *NodeDatabase) getBootstrap() string {
	x := 0
	for x != len(db.BootstrapNodeAddrs)-1 {
		if TestIP(db.BootstrapNodeAddrs[x]) {
			return db.BootstrapNodeAddrs[x]
		}
		x++
	}
	return ""
}

// NewNodeDatabase - return new node database initialized with self ID
func NewNodeDatabase(selfRef NodeID, selfAddr string) (*NodeDatabase, error) {
	readDb, err := ReadDbFromMemory(common.GetCurrentDir())

	if err != nil {
		if !strings.Contains(err.Error(), "no such file") && !strings.Contains(err.Error(), "cannot find the file") {
			return nil, err
		}
	}

	if readDb != nil {
		fmt.Println("read existing node database from mem")
		return readDb, nil
	}
	var tempArr []string
	tempArr = append(tempArr, bootStrapNode1Addr)
	return &NodeDatabase{SelfRef: selfRef, SelfAddr: selfAddr, BootstrapNodeAddrs: tempArr}, nil
}

// AddNode - add specified IP address & ID to node directory
func (db *NodeDatabase) AddNode(ip string, id NodeID) {
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
func (db *NodeDatabase) WriteDbToMemory(path string) error {
	err := common.WriteGob(path+"nodeDb.gob", db)

	if err != nil {
		fmt.Println(err)
		return err
	}

	common.ThrowSuccess("\nobject written to memory")

	return nil
}

// ReadDbFromMemory - read serialized object of specified node database from specified path
func ReadDbFromMemory(path string) (*NodeDatabase, error) {
	tempDb := new(NodeDatabase)

	err := common.ReadGob(path+"nodeDb.gob", tempDb)
	if err != nil {
		return nil, err
	}
	return tempDb, nil
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
func (db *NodeDatabase) LastPing(id NodeID) time.Time {
	nodeIndex := db.GetNodeIndex(id)
	return db.NodePingTimeDB[nodeIndex]
}

// GetNodeIndex - fetch/retrieve node index from node reference
func (db *NodeDatabase) GetNodeIndex(id NodeID) int {
	for k, v := range db.NodeRefDB {
		if id == v {
			return k
		}
	}
	return 0
}
