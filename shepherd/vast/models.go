package vast

import "time"

type VastOffer struct {
	ID        int           `json:"id"`
	GPUName   string        `json:"gpu_name"`
	GPUNumber int           `json:"num_gpus"`
	SSHHost   string        `json:"ssh_host"`
	SSHPort   int           `json:"ssh_port"`
	StartDate time.Duration `json:"start_date"`
}
