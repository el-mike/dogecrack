package vast

import (
	"github.com/el-mike/dogecrack/shepherd/provider"
)

const ProviderName = "vast"

// VastOffer - representation of a Vast.ai instance offer.
type VastOffer struct {
	ID        int    `json:"id"`
	GPUName   string `json:"gpu_name"`
	GPUNumber int    `json:"num_gpus"`
	SSHHost   string `json:"ssh_host"`
	SSHPort   int    `json:"ssh_port"`
}

type VastCreateResponse struct {
	Success    bool `json:"success"`
	InstanceId int  `json:"new_contract"`
}

// VastInstance - actual, running instance of Vast.ai machine.
type VastInstance struct {
	ID          int    `json:"id"`
	Status      string `json:"actual_status"`
	SSHHost     string `json:"ssh_host"`
	SSHPort     int    `json:"ssh_port"`
	DockerImage string `json:"image_uuid"`
}

// ProviderName - ProviderInstance implementation.
func (vi *VastInstance) ProviderName() string {
	return ProviderName
}

// ProviderId - ProviderInstance implementation.
func (vi *VastInstance) ProviderId() int {
	return vi.ID
}

// PitbullStatus - ProviderInstance implementation.
func (vi *VastInstance) PitbullStatus() provider.InstanceStatus {
	if vi.Status == "loading" {
		return provider.Starting
	}

	if vi.Status == "running" {
		return provider.Running
	}

	return provider.Finished
}
