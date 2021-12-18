package vast

import (
	"fmt"

	"github.com/el-mike/dogecrack/shepherd/pitbull"
)

// VastManager - entity responsible for managing Vast.ai machine instances.
// Please note that it depends on vast.ai CLI (vast) being added to os PATH.
type VastManager struct {
	cli *VastCLI
}

// NewVastManager - returns new VastManager instance.
func NewVastManager(apiSecret string) *VastManager {
	return &VastManager{
		cli: NewVastCLI(apiSecret),
	}
}

func (vm *VastManager) Sync() ([]pitbull.PitbullInstance, error) {
	instances, err := vm.cli.GetInstances()
	if err != nil {
		return nil, err
	}

	fmt.Print(instances)

	return nil, nil
}

func (vm *VastManager) RunInstance(fileUrl, wallet string) (pitbull.PitbullInstance, error) {
	return nil, nil
}

func (vm *VastManager) CheckInstance(instance pitbull.PitbullInstance) (pitbull.PitbullStatus, string, error) {
	return pitbull.Finished, "", nil
}
