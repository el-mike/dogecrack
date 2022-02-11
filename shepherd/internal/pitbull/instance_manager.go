package pitbull

import (
	"errors"
	"fmt"

	"github.com/el-mike/dogecrack/shepherd/internal/host"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/models"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/repositories"
)

// InstanceManager - main managing entity responsible for Pitbull instances.
type InstanceManager struct {
	hostManager        host.HostManager
	instanceRepository *repositories.InstanceRepository
}

// NewInstanceManager - returns new Shepherd instance.
func NewInstanceManager(hostManager host.HostManager) *InstanceManager {
	return &InstanceManager{
		hostManager:        hostManager,
		instanceRepository: repositories.NewInstanceRepository(),
	}
}

// GetActiveInstances - returns all active PitbullInstances.
func (ma *InstanceManager) GetActiveInstances() ([]*models.PitbullInstance, error) {
	return ma.instanceRepository.GetActiveInstances()
}

// GetInstanceById - returns a PitbullInstance with given id.
func (ma *InstanceManager) GetInstanceById(id string) (*models.PitbullInstance, error) {
	return ma.instanceRepository.GetInstanceById(id)
}

func (ma *InstanceManager) SyncInstance(id string) (*models.PitbullInstance, error) {
	pitbullInstance, err := ma.GetInstanceById(id)
	if err != nil {
		return nil, err
	}

	hostInstanceId := pitbullInstance.HostInstance.ProviderId()

	hostInstance, err := ma.hostManager.GetInstance(hostInstanceId)
	if err != nil {
		return nil, err
	}

	if hostInstance == nil {
		return nil, host.NewHostInstanceNotAvailable(hostInstanceId)
	}

	pitbullInstance.HostInstance = hostInstance

	// We want to update pitbullInstance's status and progress when host is in "running" state.
	if hostInstance.HostStatus() == host.Running {
		statusRaw, err := ma.hostManager.GetPitbullStatus(pitbullInstance.HostInstance)
		if err != nil {
			return nil, err
		}

		progressRaw, err := ma.hostManager.GetPitbullProgress(pitbullInstance.HostInstance)
		if err != nil {
			return nil, err
		}

		pitbullInstance.SetStatus(statusRaw)

		if err := pitbullInstance.SetProgress(progressRaw); err != nil {
			return nil, err
		}
	}

	if err := ma.UpdateInstance(pitbullInstance); err != nil {
		return nil, err
	}

	return pitbullInstance, nil
}

func (ma *InstanceManager) CreateInstance(passlistUrl, walletString string) (*models.PitbullInstance, error) {
	hostInstance := ma.hostManager.CreateInstance()
	pitbullInstance := models.NewPitbullInstance(hostInstance, passlistUrl, walletString)

	if err := ma.instanceRepository.CreateInstance(pitbullInstance); err != nil {
		return nil, err
	}

	return pitbullInstance, nil
}

// RunInstance - runs single pitbull instance.
func (ma *InstanceManager) RunHostForInstance(id string) (*models.PitbullInstance, error) {
	pitbullInstance, err := ma.GetInstanceById(id)
	if err != nil {
		return nil, err
	}

	if pitbullInstance.PasslistUrl == "" || pitbullInstance.WalletString == "" {
		return nil, errors.New(fmt.Sprintf("PasslistUrl or WalletString missing for instance: %s", pitbullInstance.ID.Hex()))
	}

	hostInstance, err := ma.hostManager.RunInstance()
	if err != nil {
		return nil, err
	}

	pitbullInstance.SetHost(hostInstance)
	pitbullInstance.Status = models.HostStarting

	if err := ma.instanceRepository.UpdateInstance(pitbullInstance); err != nil {
		return nil, err
	}

	return pitbullInstance, nil
}

// StopHostInstance - stops a host instance with given id.
func (ma *InstanceManager) StopHostInstance(id string) error {
	pitbullInstance, err := ma.GetInstanceById(id)
	if err != nil {
		return err
	}

	return ma.hostManager.DestroyInstance(pitbullInstance.HostInstance.ProviderId())
}

func (ma *InstanceManager) RunPitbull(id string) (*models.PitbullInstance, error) {
	pitbullInstance, err := ma.GetInstanceById(id)
	if err != nil {
		return nil, err
	}

	if pitbullInstance.HostInstance == nil ||
		pitbullInstance.PasslistUrl == "" ||
		pitbullInstance.WalletString == "" {
		return nil, errors.New(fmt.Sprintf("Instance '%s' is missing data required for running Pitbull", pitbullInstance.ID.Hex()))
	}

	if err := ma.hostManager.RunPitbull(pitbullInstance.HostInstance, pitbullInstance.PasslistUrl, pitbullInstance.WalletString); err != nil {
		return nil, err
	}

	return pitbullInstance, nil
}

// UpdateInstance - updates Pitbull instance.
func (ma *InstanceManager) UpdateInstance(pitbullInstance *models.PitbullInstance) error {
	return ma.instanceRepository.UpdateInstance(pitbullInstance)
}

// GetInstanceOutput - returns Pitbull process output for given instance.
func (ma *InstanceManager) GetInstanceOutput(pitbullInstance *models.PitbullInstance) (string, error) {
	return ma.hostManager.GetPitbullOutput(pitbullInstance.HostInstance)
}

// RunCommand - runs a command on Pitbull's host.
func (ma *InstanceManager) RunHostCommand(id, cmd string) (string, error) {
	pitbullInstance, err := ma.instanceRepository.GetInstanceById(id)
	if err != nil {
		return "", err
	}

	hostInstanceId := pitbullInstance.HostInstance.ProviderId()

	hostInstance, err := ma.hostManager.GetInstance(hostInstanceId)
	if err != nil {
		return "", err
	}

	return ma.hostManager.RunDirectCommand(hostInstance, cmd)
}
