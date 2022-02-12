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
func (im *InstanceManager) GetActiveInstances() ([]*models.PitbullInstance, error) {
	return im.instanceRepository.GetActiveInstances()
}

// GetInstanceById - returns a PitbullInstance with given id.
func (im *InstanceManager) GetInstanceById(id string) (*models.PitbullInstance, error) {
	return im.instanceRepository.GetInstanceById(id)
}

func (im *InstanceManager) SyncInstance(id string) (*models.PitbullInstance, error) {
	pitbullInstance, err := im.GetInstanceById(id)
	if err != nil {
		return nil, err
	}

	hostInstanceId := pitbullInstance.HostInstance.ProviderId()

	hostInstance, err := im.hostManager.GetInstance(hostInstanceId)
	if err != nil {
		return nil, err
	}

	if hostInstance == nil {
		return nil, host.NewHostInstanceNotAvailable(hostInstanceId)
	}

	pitbullInstance.HostInstance = hostInstance

	// We want to update pitbullInstance's status and progress when host is in "running" state.
	if hostInstance.HostStatus() == host.Running {
		statusRaw, err := im.hostManager.GetPitbullStatus(pitbullInstance.HostInstance)
		if err != nil {
			return nil, err
		}

		progressRaw, err := im.hostManager.GetPitbullProgress(pitbullInstance.HostInstance)
		if err != nil {
			return nil, err
		}

		pitbullInstance.SetStatus(statusRaw)

		if err := pitbullInstance.SetProgress(progressRaw); err != nil {
			return nil, err
		}
	}

	if err := im.UpdateInstance(pitbullInstance); err != nil {
		return nil, err
	}

	return pitbullInstance, nil
}

func (im *InstanceManager) CreateInstance(passlistUrl, walletString string) (*models.PitbullInstance, error) {
	hostInstance := im.hostManager.CreateInstance()
	pitbullInstance := models.NewPitbullInstance(hostInstance, passlistUrl, walletString)

	if err := im.instanceRepository.CreateInstance(pitbullInstance); err != nil {
		return nil, err
	}

	return pitbullInstance, nil
}

// RunInstance - runs single pitbull instance.
func (im *InstanceManager) RunHostForInstance(id string) (*models.PitbullInstance, error) {
	instance, err := im.GetInstanceById(id)
	if err != nil {
		return nil, err
	}

	if instance.PasslistUrl == "" || instance.WalletString == "" {
		return nil, errors.New(fmt.Sprintf("PasslistUrl or WalletString missing for instance: %s", instance.ID.Hex()))
	}

	hostInstance, err := im.hostManager.RunInstance()
	if err != nil {
		return nil, err
	}

	instance.SetHost(hostInstance)
	instance.Status = models.HostStarting

	if err := im.instanceRepository.UpdateInstance(instance); err != nil {
		return nil, err
	}

	return instance, nil
}

// StopHostInstance - stops a host instance with given id.
func (im *InstanceManager) StopHostInstance(id string) error {
	instance, err := im.GetInstanceById(id)
	if err != nil {
		return err
	}

	return im.hostManager.DestroyInstance(instance.HostInstance.ProviderId())
}

func (im *InstanceManager) RunPitbull(id string) (*models.PitbullInstance, error) {
	instance, err := im.GetInstanceById(id)
	if err != nil {
		return nil, err
	}

	if instance.HostInstance == nil ||
		instance.PasslistUrl == "" ||
		instance.WalletString == "" {
		return nil, errors.New(fmt.Sprintf("Instance '%s' is missing data required for running Pitbull", instance.ID.Hex()))
	}

	if err := im.hostManager.RunPitbull(instance.HostInstance, instance.PasslistUrl, instance.WalletString); err != nil {
		return nil, err
	}

	return instance, nil
}

// UpdateInstance - updates Pitbull instance.
func (im *InstanceManager) UpdateInstance(instance *models.PitbullInstance) error {
	return im.instanceRepository.UpdateInstance(instance)
}

// GetInstanceOutput - returns Pitbull process output for given instance.
func (im *InstanceManager) GetInstanceOutput(instance *models.PitbullInstance) (string, error) {
	return im.hostManager.GetPitbullOutput(instance.HostInstance)
}

// GetOrphanInstances - returns "orphan" instances.
func (im *InstanceManager) GetOrphanInstances() ([]*models.PitbullInstance, error) {
	return im.instanceRepository.GetOrphanInstances()
}

// RunCommand - runs a command on Pitbull's host.
func (im *InstanceManager) RunHostCommand(id, cmd string) (string, error) {
	instance, err := im.instanceRepository.GetInstanceById(id)
	if err != nil {
		return "", err
	}

	hostInstanceId := instance.HostInstance.ProviderId()

	hostInstance, err := im.hostManager.GetInstance(hostInstanceId)
	if err != nil {
		return "", err
	}

	return im.hostManager.RunDirectCommand(hostInstance, cmd)
}
