package fritzbox

import (
	"bytes"
	"log"
	"text/template"
)

type action string

func (a action) String() string {
	return string(a)
}

const (
	getCommonLinkProperties action = "GetCommonLinkProperties"
	getTotalBytesSent       action = "GetTotalBytesSent"
	getTotalBytesReceived   action = "GetTotalBytesReceived"
	getTotalPacketsSent     action = "GetTotalPacketsSent"
	getTotalPacketsReceived action = "GetTotalPacketsReceived"
)

type requestParams struct {
	path   string
	Action action
	Urn    string
}

var requests = map[action]requestParams{
	getCommonLinkProperties: {path: "/upnp/control/wancommonifconfig1", Action: getCommonLinkProperties, Urn: "urn:dslforum-org:service:WANCommonInterfaceConfig:1"},
	getTotalBytesSent:       {path: "/upnp/control/wancommonifconfig1", Action: getTotalBytesSent, Urn: "urn:dslforum-org:service:WANCommonInterfaceConfig:1"},
	getTotalBytesReceived:   {path: "/upnp/control/wancommonifconfig1", Action: getTotalBytesReceived, Urn: "urn:dslforum-org:service:WANCommonInterfaceConfig:1"},
	getTotalPacketsSent:     {path: "/upnp/control/wancommonifconfig1", Action: getTotalPacketsSent, Urn: "urn:dslforum-org:service:WANCommonInterfaceConfig:1"},
	getTotalPacketsReceived: {path: "/upnp/control/wancommonifconfig1", Action: getTotalPacketsReceived, Urn: "urn:dslforum-org:service:WANCommonInterfaceConfig:1"},
}

var requestTemplate = `<?xml version="1.0"?><s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body><u:{{.Action}} xmlns:u="{{.Urn}}"/></s:Body></s:Envelope>`

func buildRequestBody(params requestParams) []byte {
	t := template.Must(template.New("request").Parse(requestTemplate))
	var b bytes.Buffer
	if err := t.Execute(&b, params); err != nil {
		log.Printf("Error parsing template: %s\n", err)
	}
	return b.Bytes()
}
