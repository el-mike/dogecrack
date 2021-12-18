package vast

import "github.com/el-mike/dogecrack/shepherd/provider"

// VastManager - entity responsible for managing Vast.ai machine instances.
// Please note that it depends on vast.ai CLI (vast) being added to os PATH.
type VastManager struct {
	cli *VastCLIClient
}

// NewVastManager - returns new VastManager instance.
func NewVastManager(apiSecret string) *VastManager {
	return &VastManager{
		cli: NewVastCLI(apiSecret),
	}
}

// Sync - ProviderInstanceManager implementation.
func (vm *VastManager) Sync() ([]provider.ProviderInstance, error) {
	instances, err := vm.cli.GetInstances()
	if err != nil {
		return nil, err
	}

	providerInstances := make([]provider.ProviderInstance, len(instances))

	for _, instance := range instances {
		providerInstances = append(providerInstances, instance)
	}

	return providerInstances, nil
}

// RunInstance - ProviderInstanceManager implementation.
func (vm *VastManager) RunInstance(fileUrl, wallet string) (provider.ProviderInstance, error) {
	return nil, nil
}

// CheckInstance - ProviderInstanceManager implementation.
func (vm *VastManager) CheckInstance(instance provider.ProviderInstance) (provider.InstanceStatus, string, error) {
	return provider.Finished, "", nil
}
