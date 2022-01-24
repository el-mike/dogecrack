package vast

import (
	"errors"

	"github.com/el-mike/dogecrack/shepherd/host"
)

// VastManager - entity responsible for managing Vast.ai machine instances.
// Please note that it depends on vast.ai CLI (vast) being added to os PATH.
type VastManager struct {
	cli VastClient

	sshUser     string
	sshPassword string
	sshDir      string
}

// NewVastManager - returns new VastManager instance.
func NewVastManager(apiSecret, pitbullImage, sshUser, sshPassword, sshDir string) *VastManager {
	return &VastManager{
		// cli: NewVastCLI(apiSecret, pitbullImage),
		cli: NewVastCLIClientMock(),

		sshUser:     sshUser,
		sshPassword: sshPassword,
		sshDir:      sshDir,
	}
}

// Sync - HostManager implementation.
func (vm *VastManager) Sync() ([]host.HostInstance, error) {
	instances, err := vm.cli.GetInstances()
	if err != nil {
		return nil, err
	}

	hostInstances := make([]host.HostInstance, len(instances))

	for _, instance := range instances {
		hostInstances = append(hostInstances, instance)
	}

	return hostInstances, nil
}

// RunInstance - HostManager implementation.
func (vm *VastManager) RunInstance(fileUrl, wallet string) (host.HostInstance, error) {
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

// CheckInstance - HostManager implementation.
func (vm *VastManager) GetInstance(instanceId int) (host.HostInstance, error) {
	return vm.cli.GetInstance(instanceId)
}

// GetPitbullStatus - HostManager implementation.
func (vm *VastManager) GetPitbullStatus(instance host.HostInstance) (string, error) {
	sshClient, err := vm.getSSHClient(instance)
	if err != nil {
		return "", err
	}

	return sshClient.GetPitbullStatus()
}

// GetPitbullProgress - HostManager implementation.
func (vm *VastManager) GetPitbullProgress(instance host.HostInstance) (string, error) {
	sshClient, err := vm.getSSHClient(instance)
	if err != nil {
		return "", err
	}

	return sshClient.GetPitbullProgress()
}

// RunDirectCommand - runs given command directly on underlying host instance.
func (vm *VastManager) RunDirectCommand(instance host.HostInstance, cmd string) (string, error) {
	sshClient, err := vm.getSSHClient(instance)
	if err != nil {
		return "", err
	}

	return sshClient.run(cmd)
}

// getSSHClient - helper function for getting sshClient instance.
func (vm *VastManager) getSSHClient(instance host.HostInstance) (*VastSSHClient, error) {
	vastInstance, ok := instance.(*VastInstance)
	if !ok {
		return nil, errors.New("HostInstance is not VastInstance!")
	}

	sshClient, err := NewVastSSHClient(vm.sshUser, vm.sshPassword, vm.sshDir, vastInstance.SSHHost, vastInstance.SSHPort)
	if err != nil {
		return nil, err
	}

	return sshClient, nil
}
