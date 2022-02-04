package models

import (
	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	common.BaseModel `bson:",inline"`

	InstanceId string    `bson:"instanceId" json:"instanceId"`
	Status     JobStatus `bson:"status" json:"status"`
}

// NewPitbullJob - returns new PitbullJob instance.
func NewPitbullJob(instance *PitbullInstance) *PitbullJob {
	job := &PitbullJob{
		InstanceId: instance.ID.Hex(),
		Status:     Scheduled,
	}

	job.ID = primitive.NewObjectID()

	return job
}
