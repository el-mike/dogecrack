package vast

// VastManager - entity responsible for managing Vast.ai machine instances.
// Please note that it depends on vast.ai CLI (vast) being added to os PATH.
type VastManager struct{}

// NewVastManager - returns new VastManager instance.
func NewVastManager() *VastManager {
	return &VastManager{}
}
