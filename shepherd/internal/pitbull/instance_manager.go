package pitbull

import (
	"errors"
	"fmt"

	"github.com/el-mike/dogecrack/shepherd/internal/common/host"
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"github.com/el-mike/dogecrack/shepherd/internal/vast"
	vastmodels "github.com/el-mike/dogecrack/shepherd/internal/vast/models"
)

// InstanceManager - main managing entity responsible for Pitbull instances.
type InstanceManager struct {
	hostManager        host.HostManager
	instanceRepository *InstanceRepository
}

// NewInstanceManager - returns new Shepherd instance.
func NewInstanceManager() *InstanceManager {
	appConfig := config.GetAppConfig()

	var hostManager host.HostManager

	if appConfig.HostProvider == vastmodels.ProviderName {
		vastManager := vast.NewVastManager(appConfig.VastApiSecret, appConfig.PitbullImage, appConfig.SSHUser, appConfig.SSHPassword, appConfig.SSHDirPath, appConfig.SSHPrivateKey, appConfig.RootPath)
		hostManager = host.HostManager(vastManager)
	}

	return &InstanceManager{
		hostManager:        hostManager,
		instanceRepository: NewInstanceRepository(),
	}
}

// GetActiveInstances - returns all active PitbullInstances.
func (im *InstanceManager) GetActiveInstances() ([]*models.PitbullInstance, error) {
	return im.instanceRepository.GetActiveInstances()
}

// GetStatistics - returns PitbullInstance statistics.
func (im *InstanceManager) GetStatistics() (*models.PitbullInstancesStatistics, error) {
	return im.instanceRepository.GetStatistics()
}

// GetInstanceById - returns a PitbullInstance with given id.
func (im *InstanceManager) GetInstanceById(id string) (*models.PitbullInstance, error) {
	return im.instanceRepository.GetInstanceById(id)
}

func (im *InstanceManager) SyncInstance(id string) (*models.PitbullInstance, error) {
	instance, err := im.GetInstanceById(id)
	if err != nil {
		return nil, err
	}

	hostInstanceId := instance.HostInstance.ProviderId()

	hostInstance, err := im.hostManager.GetInstance(hostInstanceId)
	if err != nil {
		return nil, err
	}

	if hostInstance == nil {
		return nil, host.NewHostInstanceNotAvailable(hostInstanceId)
	}

	instance.HostInstance = hostInstance

	pitbull := instance.Pitbull

	// We want to update pitbullInstance's status and progress when host is in "running" state.
	if hostInstance.HostStatus() == host.Running {
		// If previous instance's status was HostStarting, we want to update it
		// to running when host instance is Running as well.
		if instance.Status == models.PitbullInstanceStatus.HostStarting {
			instance.Status = models.PitbullInstanceStatus.Running
		}

		rawStatus, err := im.hostManager.GetPitbullStatus(instance.HostInstance)
		if err != nil {
			return nil, err
		}

		rawProgress, err := im.hostManager.GetPitbullProgress(instance.HostInstance)
		if err != nil {
			return nil, err
		}

		pitbull.ParseRawStatus(rawStatus)

		if err := pitbull.ParseProgress(rawProgress); err != nil {
			return nil, err
		}
	}

	if pitbull.Done() {
		instance.CompletedAt = models.NullableTimeNow()

		if pitbull.Status == models.PitbullStatus.Success {
			instance.Status = models.PitbullInstanceStatus.Success
		} else {
			// If Pitbull process finished (not succeded), but not all passwords have been checked,
			// that means some problem occurred and instance should be marked as Failed.
			if !pitbull.AllPasswordsChecked() {
				instance.Status = models.PitbullInstanceStatus.Failed
			} else {
				instance.Status = models.PitbullInstanceStatus.Completed
			}
		}
	}

	if err := im.UpdateInstance(instance); err != nil {
		return nil, err
	}

	return instance, nil
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
	instance.Status = models.PitbullInstanceStatus.HostStarting

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

	instance.StartedAt = models.NullableTimeNow()

	if err := im.instanceRepository.UpdateInstance(instance); err != nil {
		return nil, err
	}

	return instance, nil
}

// MarkInstanceAsInterrupted - marks given instance as "Interrupted" and updates it
// in the DB.
func (im *InstanceManager) MarkInstanceAsInterrupted(instance *models.PitbullInstance) error {
	if instance == nil {
		return nil
	}

	instance.Status = models.PitbullInstanceStatus.Interrupted

	return im.UpdateInstance(instance)
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
