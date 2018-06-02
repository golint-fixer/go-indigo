package networking

import (
	"fmt"
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
