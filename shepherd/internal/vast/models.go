package vast

import (
	"github.com/el-mike/dogecrack/shepherd/internal/host"
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

// ProviderName - HostInstance implementation.
func (vi *VastInstance) ProviderName() string {
	return ProviderName
}

// ProviderId - HostInstance implementation.
func (vi *VastInstance) ProviderId() int {
	return vi.ID
}

// HostStatus - HostInstance implementation.
func (vi *VastInstance) HostStatus() host.HostStatus {
	if vi.Status == "" || vi.Status == "loading" {
		return host.Starting
	}

	if vi.Status == "running" {
		return host.Running
	}

	return host.Finished
}
