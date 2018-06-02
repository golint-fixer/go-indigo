package upnp

import (
	// "log"
	"errors"
	"net"
	"strings"
)

// GetLocalIntenetIP -
func GetLocalIntenetIP() string {
	/*
	  获得所有本机地址
	  判断能联网的ip地址
	*/

	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		panic(errors.New("不能连接网络"))
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

// GetLocalIPs -
func GetLocalIPs() ([]*net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	ips := make([]*net.IP, 0)
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		if ipnet.IP.IsLoopback() {
			continue
		}

		ips = append(ips, &ipnet.IP)
	}

	return ips, nil
}
