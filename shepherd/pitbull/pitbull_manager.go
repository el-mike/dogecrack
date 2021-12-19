package pitbull

import (
	"log"

	"github.com/el-mike/dogecrack/shepherd/models"
	"github.com/el-mike/dogecrack/shepherd/persist"
	"github.com/el-mike/dogecrack/shepherd/provider"
)

// PitbullManager - main managing entity responsible for Pitbull instances.
type PitbullManager struct {
	providerInstanceManager provider.ProviderInstanceManager
	instanceRepository      *persist.InstanceRepository
}

// NewPitbullManager - returns new Shepherd instance.
func NewPitbullManager(providerInstanceManager provider.ProviderInstanceManager) *PitbullManager {
	return &PitbullManager{
		providerInstanceManager: providerInstanceManager,
		instanceRepository:      persist.NewInstanceRepository(),
	}
}

// SyncInstances - checks and syncs given provider's instances with internal representation.
func (pm *PitbullManager) SyncInstances() error {
	instances, err := pm.providerInstanceManager.Sync()
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
	return pm.instanceRepository.GetInstanceById(id)
}

// RunInstance - runs single pitbull instance.
func (pm *PitbullManager) RunInstance(fileUrl, walletString string) (*models.PitbullInstance, error) {
	instance, err := pm.providerInstanceManager.RunInstance(fileUrl, walletString)
	if err != nil {
		return nil, err
	}

	pitbull := models.NewPitbullInstance(instance)

	if err := pm.instanceRepository.SaveInstance(pitbull); err != nil {
		return nil, err
	}

	return pitbull, nil
}
