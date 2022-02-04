package models

// JobStatus - describes the status of a single PitbullJob.
type JobStatus int8

const (
	Scheduled JobStatus = iota
	Processing
	Rescheduled
	Rejected
	Acknowledged
)

// PitbullJob - represents a single Pitbull job.
type PitbullJob struct {
	BaseModel `bson:",inline"`

	InstanceId string    `bson:"instanceId" json:"instanceId"`
	Status     JobStatus `bson:"status" json:"status"`
}

// NewPitbullJob - returns new PitbullJob instance.
func NewPitbullJob(instance *PitbullInstance) *PitbullJob {
	return &PitbullJob{
		InstanceId: instance.ID.Hex(),
		Status:     Scheduled,
	}
}
