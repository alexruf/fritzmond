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
	WANAccessTypeDSL         = "DSL"
	WANAccessTypeEthernet    = "Ethernet"
	WANAccessTypeXAVMDEFiber = "X_AVM-DE_Fiber"
	WANAccessTypeXAVMDEUMTS  = "X_AVM-DE_UMTS"
	WANAccessTypeXAVMDECable = "X_AVM-DE_Cable"
	WANAccessTypeXAVMDELTE   = "X_AVM-DE_LTE"
)
