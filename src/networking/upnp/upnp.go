package upnp

import (
	// "fmt"
	"errors"
	"log"
	"sync"
)

/*
 * 得到网关
 */

// MappingPortStruct -
type MappingPortStruct struct {
	lock         *sync.Mutex
	mappingPorts map[string][][]int
}

//添加一个端口映射记录
//只对映射进行管理
func (selfRef *MappingPortStruct) addMapping(localPort, remotePort int, protocol string) {

	selfRef.lock.Lock()
	defer selfRef.lock.Unlock()
	if selfRef.mappingPorts == nil {
		one := make([]int, 0)
		one = append(one, localPort)
		two := make([]int, 0)
		two = append(two, remotePort)
		portMapping := [][]int{one, two}
		selfRef.mappingPorts = map[string][][]int{protocol: portMapping}
		return
	}
	portMapping := selfRef.mappingPorts[protocol]
	if portMapping == nil {
		one := make([]int, 0)
		one = append(one, localPort)
		two := make([]int, 0)
		two = append(two, remotePort)
		selfRef.mappingPorts[protocol] = [][]int{one, two}
		return
	}
	one := portMapping[0]
	two := portMapping[1]
	one = append(one, localPort)
	two = append(two, remotePort)
	selfRef.mappingPorts[protocol] = [][]int{one, two}
}

//删除一个映射记录
//只对映射进行管理
func (selfRef *MappingPortStruct) delMapping(remotePort int, protocol string) {
	selfRef.lock.Lock()
	defer selfRef.lock.Unlock()
	if selfRef.mappingPorts == nil {
		return
	}
	tmp := MappingPortStruct{lock: new(sync.Mutex)}
	mappings := selfRef.mappingPorts[protocol]
	for i := 0; i < len(mappings[0]); i++ {
		if mappings[1][i] == remotePort {
			//要删除的映射
			break
		}
		tmp.addMapping(mappings[0][i], mappings[1][i], protocol)
	}
	selfRef.mappingPorts = tmp.mappingPorts
}

// GetAllMapping -
func (selfRef *MappingPortStruct) GetAllMapping() map[string][][]int {
	return selfRef.mappingPorts
}

// Upnp -
type Upnp struct {
	Active             bool              //这个upnp协议是否可用
	LocalHost          string            //本机ip地址
	GatewayInsideIP    string            //局域网网关ip
	GatewayOutsideIP   string            //网关公网ip
	OutsideMappingPort map[string]int    //映射外部端口
	InsideMappingPort  map[string]int    //映射本机端口
	Gateway            *Gateway          //网关信息
	CtrlURL            string            //控制请求url
	MappingPort        MappingPortStruct //已经添加了的映射 {"TCP":[1990],"UDP":[1991]}
}

// SearchGateway -
func (selfRef *Upnp) SearchGateway() (err error) {
	defer func(err error) {
		if errTemp := recover(); errTemp != nil {
			log.Println("upnp module error", errTemp)
			err = errTemp.(error)
		}
	}(err)

	if selfRef.LocalHost == "" {
		selfRef.MappingPort = MappingPortStruct{
			lock: new(sync.Mutex),
			// mappingPorts: map[string][][]int{},
		}
		selfRef.LocalHost = GetLocalIntenetIP()
	}
	searchGateway := SearchGateway{upnp: selfRef}
	if searchGateway.Send() {
		return nil
	}
	return errors.New("no gateway device found")
}

func (selfRef *Upnp) deviceStatus() {

}

//查看设备描述，得到控制请求url
func (selfRef *Upnp) deviceDesc() (err error) {
	if selfRef.GatewayInsideIP == "" {
		if err := selfRef.SearchGateway(); err != nil {
			return err
		}
	}
	device := DeviceDesc{upnp: selfRef}
	device.Send()
	selfRef.Active = true
	// log.Println("获得控制请求url:", selfRef.CtrlUrl)
	return
}

// ExternalIPAddr -
func (selfRef *Upnp) ExternalIPAddr() (err error) {
	if selfRef.CtrlURL == "" {
		if err := selfRef.deviceDesc(); err != nil {
			return err
		}
	}
	eia := ExternalIPAddress{upnp: selfRef}
	eia.Send()
	return nil
	// log.Println("获得公网ip地址为：", selfRef.GatewayOutsideIP)
}

// AddPortMapping -
func (selfRef *Upnp) AddPortMapping(localPort, remotePort int, protocol string) (err error) {
	defer func(err error) {
		if errTemp := recover(); errTemp != nil {
			log.Println("upnp module error", errTemp)
			err = errTemp.(error)
		}
	}(err)
	if selfRef.GatewayOutsideIP == "" {
		if err := selfRef.ExternalIPAddr(); err != nil {
			return err
		}
	}
	addPort := AddPortMapping{upnp: selfRef}
	if issuccess := addPort.Send(localPort, remotePort, protocol); issuccess {
		selfRef.MappingPort.addMapping(localPort, remotePort, protocol)
		// log.Println("添加一个端口映射：protocol:", protocol, "local:", localPort, "remote:", remotePort)
		return nil
	}
	selfRef.Active = false
	return errors.New("failed to add port mapping")
}

// DelPortMapping -
func (selfRef *Upnp) DelPortMapping(remotePort int, protocol string) bool {
	delMapping := DelPortMapping{upnp: selfRef}
	issuccess := delMapping.Send(remotePort, protocol)
	if issuccess {
		selfRef.MappingPort.delMapping(remotePort, protocol)
		log.Println("removed port mapping： remote:", remotePort)
	}
	return issuccess
}

// Reclaim -
func (selfRef *Upnp) Reclaim() {
	mappings := selfRef.MappingPort.GetAllMapping()
	tcpMapping, ok := mappings["TCP"]
	if ok {
		for i := 0; i < len(tcpMapping[0]); i++ {
			selfRef.DelPortMapping(tcpMapping[1][i], "TCP")
		}
	}
	udpMapping, ok := mappings["UDP"]
	if ok {
		for i := 0; i < len(udpMapping[0]); i++ {
			selfRef.DelPortMapping(udpMapping[0][i], "UDP")
		}
	}
}

// GetAllMapping -
func (selfRef *Upnp) GetAllMapping() map[string][][]int {
	return selfRef.MappingPort.GetAllMapping()
}
