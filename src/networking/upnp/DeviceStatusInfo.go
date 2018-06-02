package upnp

import (
	// "log"
	// "io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// SearchGatewayReq -
type SearchGatewayReq struct {
	host       string
	resultBody string
	ctrlURL    string
	upnp       *Upnp
}

// Send -
func (selfRef SearchGatewayReq) Send() {
	// request := selfRef.BuildRequest()
}

// BuildRequest -
func (selfRef SearchGatewayReq) BuildRequest() *http.Request {
	//请求头
	header := http.Header{}
	header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	header.Set("SOAPAction", `"urn:schemas-upnp-org:service:WANIPConnection:1#GetStatusInfo"`)
	header.Set("Content-Type", "text/xml")
	header.Set("Connection", "Close")
	header.Set("Content-Length", "")
	//请求体
	body := Node{Name: "SOAP-ENV:Envelope",
		Attr: map[string]string{"xmlns:SOAP-ENV": `"http://schemas.xmlsoap.org/soap/envelope/"`,
			"SOAP-ENV:encodingStyle": `"http://schemas.xmlsoap.org/soap/encoding/"`}}
	childOne := Node{Name: `SOAP-ENV:Body`}
	childTwo := Node{Name: `m:GetStatusInfo`,
		Attr: map[string]string{"xmlns:m": `"urn:schemas-upnp-org:service:WANIPConnection:1"`}}
	childOne.AddChild(childTwo)
	body.AddChild(childOne)
	bodyStr := body.BuildXML()
	//请求
	request, _ := http.NewRequest("POST", "http://"+selfRef.upnp.Gateway.Host+selfRef.upnp.CtrlURL,
		strings.NewReader(bodyStr))
	request.Header = header
	request.Header.Set("Content-Length", strconv.Itoa(len([]byte(bodyStr))))
	return request
}
