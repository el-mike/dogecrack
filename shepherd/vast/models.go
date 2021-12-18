package vast

import (
	"time"

	"github.com/el-mike/dogecrack/shepherd/pitbull"
)

type VastOffer struct {
	ID        int           `json:"id"`
	GPUName   string        `json:"gpu_name"`
	GPUNumber int           `json:"num_gpus"`
	SSHHost   string        `json:"ssh_host"`
	SSHPort   int           `json:"ssh_port"`
	StartDate time.Duration `json:"start_date"`
}

type VastInstance struct {
	pitbullId string
	Status    string `json:"intended_status"`
	VastId    string `json:"id"`
}

// PitbullId - PitbullInstance implementation.
func (vi *VastInstance) PitbullId() string {
	return vi.pitbullId
}

// ProviderId - PitbullInstance implementation.
func (vi *VastInstance) ProviderId() string {
	return vi.VastId
}

// PitbullStatus - PitbullInstance implementation.
func (vi *VastInstance) PitbullStatus() pitbull.PitbullStatus {
	if vi.Status == "starting" {
		return pitbull.Starting
	}

	if vi.Status == "running" {
		return pitbull.Running
	}

	return pitbull.Finished
}
