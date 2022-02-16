package models

// Offer - representation of a Vast.ai instance offer.
type Offer struct {
	ID        int    `json:"id"`
	GPUName   string `json:"gpu_name"`
	GPUNumber int    `json:"num_gpus"`
	SSHHost   string `json:"ssh_host"`
	SSHPort   int    `json:"ssh_port"`
}
