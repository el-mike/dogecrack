package pitbull

import (
	"fmt"
	"log"

	"github.com/el-mike/dogecrack/shepherd/common"
)

// PitbullManager - main managing entity responsible for Pitbull instances.
type PitbullManager struct {
	providerInstanceManager common.ProviderInstanceManager
}

// NewPitbullManager - returns new Shepherd instance.
func NewPitbullManager(providerInstanceManager common.ProviderInstanceManager) *PitbullManager {
	return &PitbullManager{
		providerInstanceManager: providerInstanceManager,
	}
}

func (pm *PitbullManager) SyncInstances() error {
	instances, err := pm.providerInstanceManager.Sync()
	if err != nil {
		log.Fatal("Instances sync failed!")

		return err
	}

	fmt.Print(instances)

	return nil
}

func (pm *PitbullManager) RunInstance(filUrl, walletString string) (*common.PitbullInstance, error) {
	return nil, nil
}
