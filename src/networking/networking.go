package networking

import (
	"fmt"
	"indo-go/src/common"
	"indo-go/src/core/types"
	"indo-go/src/networking/upnp"
	"net"
)

// AddPortMapping - add port mapping on specified port
func AddPortMapping(port int) {
	mapping := new(upnp.Upnp)
	if err := mapping.AddPortMapping(port, port, "TCP"); err == nil {
		fmt.Println("port mapping added")
	} else {
		fmt.Printf("port mapping failed with err %s", err)
	}
}

// Relay - push localized or received transaction to further node.
func Relay(Tx *types.Transaction) {

}

func (conn *Connection) attempt() {
	connec, err := net.Dial("tcp", conn.DestNodeAddr+":3000") // Connect to peer addr
	connec.Write(conn.Data)                                   // Write specified data
	if err != nil {
		fmt.Println(err)
	} else {
		conn.AddEvent("started")
	}
}

func newConnection(initAddr string, destAddr string, connType ConnectionType, data []byte) *Connection {
	if common.StringInSlice(string(connType), ConnectionTypes) {
		conn := Connection{InitNodeAddr: initAddr, DestNodeAddr: destAddr, Type: connType, Data: data}
		return &conn
	}
	common.ThrowWarning("connection type not valid")
	return nil
}

//https://systembash.com/a-simple-go-tcp-server-and-tcp-client/
