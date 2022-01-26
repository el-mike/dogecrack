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
func NewVastManager(apiSecret, pitbullImage, sshUser, sshPassword, sshDir, rootDir string) *VastManager {
	return &VastManager{
		cli: NewVastCLI(apiSecret, pitbullImage),
		// cli: NewVastCLIClientMock(rootDir),

		sshUser:     sshUser,
		sshPassword: sshPassword,
		sshDir:      sshDir,
	}
}

// CreateInstance - HostManager implementation.
func (vm *VastManager) CreateInstance() host.HostInstance {
	return &VastInstance{}
}

// RunInstance - HostManager implementation.
func (vm *VastManager) RunInstance() (host.HostInstance, error) {
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

// DestroyInstance - stops given instance.
func (vm *VastManager) DestroyInstance(instanceId int) error {
	return vm.cli.DestroyInstance(instanceId)
}

func (vm *VastManager) RunPitbull(instance host.HostInstance, passlistUrl, walletString string) error {
	sshClient, err := vm.getSSHClient(instance)
	if err != nil {
		return err
	}

	_, err = sshClient.RunPitbull(passlistUrl, walletString)
	if err != nil {
		return err
	}

	return nil
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

func (vm *VastManager) GetPitbullOutput(instance host.HostInstance) (string, error) {
	sshClient, err := vm.getSSHClient(instance)
	if err != nil {
		return "", err
	}

	return sshClient.GetPitbullOutput()
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
