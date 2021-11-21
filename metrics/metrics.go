package metrics

import "github.com/alexruf/fritzmond/fritzbox"

const (
	MetricPhysicalLinkStatus         = "physicalLinkStatus"
	MetricLayer1DownstreamMaxBitRate = "layer1DownstreamMaxBitRate"
	MetricLayer1UpstreamMaxBitRate   = "layer1UpstreamMaxBitRate"
	MetricTotalBytesSent             = "totalBytesSent"
	MetricTotalBytesReceived         = "totalBytesReceived"
	MetricTotalPacketsSent           = "totalPacketsSent"
	MetricTotalPacketsReceived       = "totalPacketsReceived"
)

// ConvertPhysicalLinkStatus converts the physical link status string to a float64 metric value.
// 1 = Up;
// 1.5 = Initializing;
// 0 = Down;
// -1 = Unavailable;
func ConvertPhysicalLinkStatus(status fritzbox.PhysicalLinkStatus) float64 {
	switch status {
	case fritzbox.PhysicalLinkStatusUp:
		return 1
	case fritzbox.PhysicalLinkStatusInitializing:
		return 1.5
	case fritzbox.PhysicalLinkStatusDown:
		return 0
	case fritzbox.PhysicalLinkStatusUnavailable:
		fallthrough
	default:
		break
	}
	return -1
}
