package vast

import (
	"errors"

	"github.com/el-mike/dogecrack/shepherd/internal/common/host"
	"github.com/el-mike/dogecrack/shepherd/internal/vast/models"
)

// VastManager - entity responsible for managing Vast.ai machine instances.
// Please note that it depends on vast.ai CLI (vast) being added to os PATH.
type VastManager struct {
	cli VastClient

	sshUser       string
	sshPassword   string
	sshDir        string
	sshPrivateKey string
}

// NewVastManager - returns new VastManager instance.
func NewVastManager(apiSecret, pitbullImage, sshUser, sshPassword, sshDir, sshPrivateKey, rootDir string) *VastManager {
	return &VastManager{
		// cli: NewVastCLI(apiSecret, pitbullImage),
		cli: NewVastCLIClientMock(rootDir),

		sshUser:       sshUser,
		sshPassword:   sshPassword,
		sshDir:        sshDir,
		sshPrivateKey: sshPrivateKey,
	}
}

// CreateInstance - HostManager implementation.
func (vm *VastManager) CreateInstance() host.HostInstance {
	return &models.Instance{}
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

// GetInstance - HostManager implementation.
func (vm *VastManager) GetInstance(instanceId int) (host.HostInstance, error) {
	return vm.cli.GetInstance(instanceId)
}

// DestroyInstance - stops given instance.
func (vm *VastManager) DestroyInstance(instanceId int) error {
	return vm.cli.DestroyInstance(instanceId)
}

// RunPitbullForPasslist - runs Pitbull with passwords stored under given passlistUrl.
func (vm *VastManager) RunPitbullForPasslist(instance host.HostInstance, walletString, passlistUrl string) error {
	sshClient, err := vm.getSSHClient(instance)
	if err != nil {
		return err
	}

	_, err = sshClient.RunPitbullForPasslist(passlistUrl, walletString)
	if err != nil {
		return err
	}

	return nil
}

// RunPitbullForTokenlist - runs Pitbull with provided tokenlist.
func (vm *VastManager) RunPitbullForTokenlist(instance host.HostInstance, walletString, tokenlist string) error {
	sshClient, err := vm.getSSHClient(instance)
	if err != nil {
		return err
	}

	_, err = sshClient.RunPitbullForTokenlist(tokenlist, walletString)
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
	vastInstance, ok := instance.(*models.Instance)
	if !ok {
		return nil, errors.New("HostInstance is not VastInstance!")
	}

	sshClient, err := NewVastSSHClient(vm.sshUser, vm.sshPassword, vm.sshDir, vm.sshPrivateKey, vastInstance.SSHHost, vastInstance.SSHPort)
	if err != nil {
		return nil, err
	}

	return sshClient, nil
}
