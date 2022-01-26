package pitbull

import (
	"errors"
	"fmt"

	"github.com/el-mike/dogecrack/shepherd/host"
	"github.com/el-mike/dogecrack/shepherd/models"
	"github.com/el-mike/dogecrack/shepherd/persist"
)

// PitbullManager - main managing entity responsible for Pitbull instances.
type PitbullManager struct {
	hostManager               host.HostManager
	pitbullInstanceRepository *persist.PitbullInstanceRepository
}

// NewPitbullManager - returns new Shepherd instance.
func NewPitbullManager(hostManager host.HostManager) *PitbullManager {
	return &PitbullManager{
		hostManager:               hostManager,
		pitbullInstanceRepository: persist.NewPitbullInstanceRepository(),
	}
}

// GetActiveInstances - returns all active PitbullInstances.
func (pm *PitbullManager) GetActiveInstances() ([]*models.PitbullInstance, error) {
	return pm.pitbullInstanceRepository.GetActiveInstances()
}

// GetInstanceById - returns a PitbullInstance with given id.
func (pm *PitbullManager) GetInstanceById(id string) (*models.PitbullInstance, error) {
	return pm.pitbullInstanceRepository.GetInstanceById(id)
}

func (pm *PitbullManager) SyncInstance(id string) (*models.PitbullInstance, error) {
	pitbullInstance, err := pm.GetInstanceById(id)
	if err != nil {
		return nil, err
	}

	hostInstanceId := pitbullInstance.HostInstance.ProviderId()

	hostInstance, err := pm.hostManager.GetInstance(hostInstanceId)
	if err != nil {
		// If given instance could not be found on given provider's side, that means
		// it has been terminated - therefore, we want to mark related PitbullInstance
		// as Finished.
		if _, ok := err.(*host.HostInstanceNotFound); ok {
			pitbullInstance.Status = models.Finished
		} else {
			return nil, err
		}
		// Otherwise, we override current ProviderInstance with new data.
	} else {
		pitbullInstance.HostInstance = hostInstance
	}

	// We want to update pitbullInstance's status and progress when host is in "running" state.
	if hostInstance.HostStatus() == host.Running {
		statusRaw, err := pm.hostManager.GetPitbullStatus(pitbullInstance.HostInstance)
		if err != nil {
			return nil, err
		}

		progressRaw, err := pm.hostManager.GetPitbullProgress(pitbullInstance.HostInstance)
		if err != nil {
			return nil, err
		}

		pitbullInstance.SetStatus(statusRaw)

		if err := pitbullInstance.SetProgress(progressRaw); err != nil {
			return nil, err
		}
	}

	if err := pm.UpdateInstance(pitbullInstance); err != nil {
		return nil, err
	}

	return pitbullInstance, nil
}

func (pm *PitbullManager) CreateInstance(passlistUrl, walletString string) (*models.PitbullInstance, error) {
	hostInstance := pm.hostManager.CreateInstance()
	pitbullInstance := models.NewPitbullInstance(hostInstance, passlistUrl, walletString)

	if err := pm.pitbullInstanceRepository.CreateInstance(pitbullInstance); err != nil {
		return nil, err
	}

	return pitbullInstance, nil
}

// RunInstance - runs single pitbull instance.
func (pm *PitbullManager) RunHostForInstance(id string) (*models.PitbullInstance, error) {
	pitbullInstance, err := pm.GetInstanceById(id)
	if err != nil {
		return nil, err
	}

	if pitbullInstance.PasslistUrl == "" || pitbullInstance.WalletString == "" {
		return nil, errors.New(fmt.Sprintf("PasslistUrl or WalletString missing for instance: %s", pitbullInstance.ID.Hex()))
	}

	hostInstance, err := pm.hostManager.RunInstance()
	if err != nil {
		return nil, err
	}

	pitbullInstance.SetHost(hostInstance)
	pitbullInstance.Status = models.Starting

	if err := pm.pitbullInstanceRepository.UpdateInstance(pitbullInstance); err != nil {
		return nil, err
	}

	return pitbullInstance, nil
}

// StopHostInstance - stops a host instance with given id.
func (pm *PitbullManager) StopHostInstance(id string) error {
	pitbullInstance, err := pm.GetInstanceById(id)
	if err != nil {
		return err
	}

	return pm.hostManager.DestroyInstance(pitbullInstance.HostInstance.ProviderId())
}

func (pm *PitbullManager) RunPitbull(id string) (*models.PitbullInstance, error) {
	pitbullInstance, err := pm.GetInstanceById(id)
	if err != nil {
		return nil, err
	}

	if pitbullInstance.HostInstance == nil ||
		pitbullInstance.PasslistUrl == "" ||
		pitbullInstance.WalletString == "" {
		return nil, errors.New(fmt.Sprintf("Instance '%s' is missing data required for running Pitbull", pitbullInstance.ID.Hex()))
	}

	if err := pm.hostManager.RunPitbull(pitbullInstance.HostInstance, pitbullInstance.PasslistUrl, pitbullInstance.WalletString); err != nil {
		return nil, err
	}

	return pitbullInstance, nil
}

// UpdateInstance - updates Pitbull instance.
func (pm *PitbullManager) UpdateInstance(pitbullInstance *models.PitbullInstance) error {
	return pm.pitbullInstanceRepository.UpdateInstance(pitbullInstance)
}

// GetInstanceOutput - returns Pitbull process output for given instance.
func (pm *PitbullManager) GetInstanceOutput(pitbullInstance *models.PitbullInstance) (string, error) {
	return pm.hostManager.GetPitbullOutput(pitbullInstance.HostInstance)
}

// RunCommand - runs a command on Pitbull's host.
func (pm *PitbullManager) RunHostCommand(id, cmd string) (string, error) {
	pitbullInstance, err := pm.pitbullInstanceRepository.GetInstanceById(id)
	if err != nil {
		return "", err
	}

	hostInstanceId := pitbullInstance.HostInstance.ProviderId()

	hostInstance, err := pm.hostManager.GetInstance(hostInstanceId)
	if err != nil {
		return "", err
	}

	return pm.hostManager.RunDirectCommand(hostInstance, cmd)
}
