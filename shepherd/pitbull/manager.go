package pitbull

import (
	"log"

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

// SyncInstances - checks and syncs given provider's instances with internal representation.
func (pm *PitbullManager) SyncInstances() error {
	instances, err := pm.hostManager.Sync()
	if err != nil {
		log.Fatal("Instances sync failed!")

		return err
	}

	var pitbulls []*models.PitbullInstance

	for _, instance := range instances {
		pitbulls = append(pitbulls, models.NewPitbullInstance(instance))
	}

	return nil
}

// GetActiveInstances - returns all active PitbullInstances.
func (pm *PitbullManager) GetActiveInstances() ([]*models.PitbullInstance, error) {
	return pm.pitbullInstanceRepository.GetActiveInstances()
}

// GetInstanceById - returns a PitbullInstance with given id.
func (pm *PitbullManager) GetInstanceById(id string) (*models.PitbullInstance, error) {
	pitbullInstance, err := pm.pitbullInstanceRepository.GetInstanceById(id)
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

	if err := pm.UpdateInstance(pitbullInstance); err != nil {
		return nil, err
	}

	return pitbullInstance, nil
}

// RunInstance - runs single pitbull instance.
func (pm *PitbullManager) RunInstance(fileUrl, walletString string) (*models.PitbullInstance, error) {
	hostInstance, err := pm.hostManager.RunInstance(fileUrl, walletString)
	if err != nil {
		return nil, err
	}

	pitbullInstance := models.NewPitbullInstance(hostInstance)

	if err := pm.pitbullInstanceRepository.CreateInstance(pitbullInstance); err != nil {
		return nil, err
	}

	return pitbullInstance, nil
}

// UpdateInstance - updates Pitbull instance status and progress.
func (pm *PitbullManager) UpdateInstance(pitbullInstance *models.PitbullInstance) error {
	statusRaw, err := pm.hostManager.GetPitbullStatus(pitbullInstance.HostInstance)
	if err != nil {
		return err
	}

	progressRaw, err := pm.hostManager.GetPitbullProgress(pitbullInstance.HostInstance)
	if err != nil {
		return err
	}

	pitbullInstance.SetStatus(statusRaw)

	if err := pitbullInstance.SetProgress(progressRaw); err != nil {
		return err
	}

	if err := pm.pitbullInstanceRepository.UpdateInstance(pitbullInstance); err != nil {
		return err
	}

	return nil
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
