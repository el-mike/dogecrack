package vast

import (
	"github.com/el-mike/dogecrack/shepherd/provider"
)

// VastManager - entity responsible for managing Vast.ai machine instances.
// Please note that it depends on vast.ai CLI (vast) being added to os PATH.
type VastManager struct {
	cli *VastCLIClient
}

// NewVastManager - returns new VastManager instance.
func NewVastManager(apiSecret, pitbullImage string) *VastManager {
	return &VastManager{
		cli: NewVastCLI(apiSecret, pitbullImage),
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
	offer, err := vm.cli.GetOfferByCriteria(CheapOfferFilter)
	if err != nil {
		return nil, err
	}

	createResponse, err := vm.cli.StartInstance(offer.ID)
	if err != nil {
		return nil, err
	}

	instance, err := vm.cli.GetInstance(createResponse.InstanceId)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

// CheckInstance - ProviderInstanceManager implementation.
func (vm *VastManager) GetInstance(instanceId int) (provider.ProviderInstance, error) {
	return vm.cli.GetInstance(instanceId)
}
