package models

// PitbullJob - represents a single Pitbull job.
type PitbullJob struct {
	InstanceId string `json:"instanceId"`
}

// NewPitbullJob - returns new PitbullJob instance.
func NewPitbullJob(instance *PitbullInstance) *PitbullJob {
	return &PitbullJob{
		InstanceId: instance.ID.Hex(),
	}
}
