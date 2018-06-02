package networking

import (
	"fmt"
	"indo-go/src/common"
	"indo-go/src/core/types"
	"indo-go/src/networking/upnp"
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

func newConnection(initAddr string, destAddr string, connType ConnectionType) *Connection {
	if common.StringInSlice(string(connType), ConnectionTypes) {
		return &Connection{InitNodeAddr: initAddr, DestNodeAddr: destAddr, Type: connType}
	}
	common.ThrowWarning("connection type not valid")
	return nil
}
