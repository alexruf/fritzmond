package fritzbox

type PhysicalLinkStatus string

func (p PhysicalLinkStatus) String() string {
	return string(p)
}

const (
	PhysicalLinkStatusUnavailable  PhysicalLinkStatus = "Unavailable"
	PhysicalLinkStatusDown         PhysicalLinkStatus = "Down"
	PhysicalLinkStatusInitializing PhysicalLinkStatus = "Initializing"
	PhysicalLinkStatusUp           PhysicalLinkStatus = "Up"
)

type WANAccessType string

func (w WANAccessType) String() string {
	return string(w)
}

const (
	WANAccessTypeDSL         WANAccessType = "DSL"
	WANAccessTypeEthernet    WANAccessType = "Ethernet"
	WANAccessTypeXAVMDEFiber WANAccessType = "X_AVM-DE_Fiber"
	WANAccessTypeXAVMDEUMTS  WANAccessType = "X_AVM-DE_UMTS"
	WANAccessTypeXAVMDECable WANAccessType = "X_AVM-DE_Cable"
	WANAccessTypeXAVMDELTE   WANAccessType = "X_AVM-DE_LTE"
)
