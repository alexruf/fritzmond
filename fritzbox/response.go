package fritzbox

import "encoding/xml"

type soapEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    soapBody
}

type soapBody struct {
	XMLName              xml.Name              `xml:"Body"`
	CommonLinkProperties *CommonLinkProperties `xml:"urn:dslforum-org:service:WANCommonInterfaceConfig:1 GetCommonLinkPropertiesResponse"`
	TotalBytesSent       *TotalBytesSent       `xml:"urn:dslforum-org:service:WANCommonInterfaceConfig:1 GetTotalBytesSentResponse"`
	TotalBytesReceived   *TotalBytesReceived   `xml:"urn:dslforum-org:service:WANCommonInterfaceConfig:1 GetTotalBytesReceivedResponse"`
	TotalPacketsSent     *TotalPacketsSent     `xml:"urn:dslforum-org:service:WANCommonInterfaceConfig:1 GetTotalPacketsSentResponse"`
	TotalPacketsReceived *TotalPacketsReceived `xml:"urn:dslforum-org:service:WANCommonInterfaceConfig:1 GetTotalPacketsReceivedResponse"`
}

type CommonLinkProperties struct {
	NewWANAccessType              WANAccessType      `xml:"NewWANAccessType"`
	NewLayer1UpstreamMaxBitRate   uint               `xml:"NewLayer1UpstreamMaxBitRate"`
	NewLayer1DownstreamMaxBitRate uint               `xml:"NewLayer1DownstreamMaxBitRate"`
	NewPhysicalLinkStatus         PhysicalLinkStatus `xml:"NewPhysicalLinkStatus"`
}

type TotalBytesSent struct {
	NewTotalBytesSent uint `xml:"NewTotalBytesSent"`
}

type TotalBytesReceived struct {
	NewTotalBytesReceived uint `xml:"NewTotalBytesReceived"`
}

type TotalPacketsSent struct {
	NewTotalPacketsSent uint `xml:"NewTotalPacketsSent"`
}

type TotalPacketsReceived struct {
	NewTotalPacketsReceived uint `xml:"NewTotalPacketsReceived"`
}
