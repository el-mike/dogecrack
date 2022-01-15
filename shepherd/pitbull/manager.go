package pitbull

import (
	"log"

	"github.com/el-mike/dogecrack/shepherd/host"
	"github.com/el-mike/dogecrack/shepherd/models"
	"github.com/el-mike/dogecrack/shepherd/persist"
)

// Host i HostManager zamiast HostInstance i HostInstanceManager

// PitbullManager - main managing entity responsible for Pitbull instances.
type PitbullManager struct {
	hostManager        host.HostManager
	instanceRepository *persist.PitbullInstanceRepository
}

// NewPitbullManager - returns new Shepherd instance.
func NewPitbullManager(hostManager host.HostManager) *PitbullManager {
	return &PitbullManager{
		hostManager:        hostManager,
		instanceRepository: persist.NewPitbullInstanceRepository(),
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
	return pm.instanceRepository.GetActiveInstances()
}

// GetInstanceById - returns a PitbullInstance with given id.
func (pm *PitbullManager) GetInstanceById(id string) (*models.PitbullInstance, error) {
	instance, err := pm.instanceRepository.GetInstanceById(id)
	if err != nil {
		return nil, err
	}

	providerInstanceId := instance.HostInstance.ProviderId()

	providerInstance, err := pm.hostManager.GetInstance(providerInstanceId)
	if err != nil {
		// If given instance could not be found on given provider's side, that means
		// it has been terminated - therefore, we want to mark related PitbullInstance
		// as Finished.
		if _, ok := err.(*host.HostInstanceNotFound); ok {
			instance.Status = models.Finished
		} else {
			return nil, err
		}
		// Otherwise, we override current ProviderInstance with new data.
	} else {
		instance.HostInstance = providerInstance
	}

	if err := pm.instanceRepository.UpdateInstance(instance); err != nil {
		return nil, err
	}

	return instance, nil
}

// RunInstance - runs single pitbull instance.
func (pm *PitbullManager) RunInstance(fileUrl, walletString string) (*models.PitbullInstance, error) {
	instance, err := pm.hostManager.RunInstance(fileUrl, walletString)
	if err != nil {
		return nil, err
	}

	pitbull := models.NewPitbullInstance(instance)

	if err := pm.instanceRepository.CreateInstance(pitbull); err != nil {
		return nil, err
	}

	return pitbull, nil
}
