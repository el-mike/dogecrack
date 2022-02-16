package models

import "github.com/el-mike/dogecrack/shepherd/internal/common/host"

// Instance - actual, running instance of Vast.ai machine.
type Instance struct {
	ID          int    `json:"id"`
	Status      string `json:"actual_status"`
	SSHHost     string `json:"ssh_host"`
	SSHPort     int    `json:"ssh_port"`
	DockerImage string `json:"image_uuid"`
	// DPHTotal - cost per hour in USD.
	DPHTotal float32 `json:"dph_total"`
	// DLPerf - perf score, as defined by vast.ai.
	DLPerf float32 `json:"dlperf"`
	// DLPerfPerDPHTotal - overall proportion of DLPerf to price.
	DLPerfPerDPHTotal float32 `json:"dlperf_per_dphtotal"`
	GPUName           string  `json:"gpu_name"`
	// GPUNum - the quantity of GPUs.
	GPUsNum int `json:"num_gpus"`
}

// ProviderName - HostInstance implementation.
func (vi *Instance) ProviderName() string {
	return ProviderName
}

// ProviderId - HostInstance implementation.
func (vi *Instance) ProviderId() int {
	return vi.ID
}

// HostStatus - HostInstance implementation.
func (vi *Instance) HostStatus() host.HostStatus {
	if vi == nil || vi.Status == "" {
		return host.Unknown
	}

	if vi.Status == "loading" {
		return host.Starting
	}

	if vi.Status == "running" {
		return host.Running
	}

	return host.Destroyed
}
