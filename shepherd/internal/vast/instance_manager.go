package vast

import (
	"errors"
	"github.com/el-mike/dogecrack/shepherd/internal/common/pitbull_client"
	"os"

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

	searchCriteriaProvider SearchCriteriaProvider
}

// NewVastManager - returns new VastManager instance.
func NewVastManager(
	apiSecret,
	pitbullImage,
	sshUser,
	sshPassword,
	sshDir,
	sshPrivateKey,
	rootDir string,
	searchCriteriaProvider SearchCriteriaProvider,
) *VastManager {
	var cliClient VastClient

	if os.Getenv("USE_MOCK_VAST_CLI") == "true" {
		cliClient = NewVastCLIClientMock(rootDir)
	} else {
		cliClient = NewVastCLI(apiSecret, pitbullImage)
	}

	return &VastManager{
		cli:           cliClient,
		sshUser:       sshUser,
		sshPassword:   sshPassword,
		sshDir:        sshDir,
		sshPrivateKey: sshPrivateKey,

		searchCriteriaProvider: searchCriteriaProvider,
	}
}

// CreateInstance - HostManager implementation.
func (vm *VastManager) CreateInstance() host.HostInstance {
	return &models.Instance{}
}

// RunInstance - HostManager implementation.
func (vm *VastManager) RunInstance() (host.HostInstance, error) {
	searchCriteria := vm.searchCriteriaProvider.GetSearchCriteria()

	offer, err := vm.cli.GetOfferByCriteria(searchCriteria)
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
func (vm *VastManager) RunPitbullForPasslist(
	instance host.HostInstance,
	walletString,
	passlistUrl string,
	skipCount,
	minLength,
	maxLength int64,
) error {
	sshClient, err := vm.getSSHClient(instance)
	if err != nil {
		return err
	}

	_, err = sshClient.RunPitbullForPasslist(walletString, passlistUrl, skipCount, minLength, maxLength)
	if err != nil {
		return err
	}

	return nil
}

// RunPitbullForTokenlist - runs Pitbull with provided tokenlist.
func (vm *VastManager) RunPitbullForTokenlist(
	instance host.HostInstance,
	walletString,
	tokenlist string,
	skipCount,
	minLength,
	maxLength int64,
) error {
	sshClient, err := vm.getSSHClient(instance)
	if err != nil {
		return err
	}

	_, err = sshClient.RunPitbullForTokenlist(walletString, tokenlist, skipCount, minLength, maxLength)
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

	return sshClient.Run(cmd)
}

// getSSHClient - helper function for getting sshClient instance.
func (vm *VastManager) getSSHClient(instance host.HostInstance) (*pitbull_client.PitbullSSHClient, error) {
	vastInstance, ok := instance.(*models.Instance)
	if !ok {
		return nil, errors.New("HostInstance is not VastInstance!")
	}

	sshClient, err := pitbull_client.NewPitbullSSHClient(vm.sshUser, vm.sshPassword, vm.sshDir, vm.sshPrivateKey, vastInstance.SSHHost, vastInstance.SSHPort)
	if err != nil {
		return nil, err
	}

	return sshClient, nil
}
