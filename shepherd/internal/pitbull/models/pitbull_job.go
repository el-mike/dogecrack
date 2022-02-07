package models

import (
	"time"

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

var jobStatusNames = map[JobStatus]string{
	Scheduled:    "SCHEDULED",
	Processing:   "PROCESSING",
	Rescheduled:  "RESCHEDULED",
	Rejected:     "REJECTED",
	Acknowledged: "ACKNOWLEDGED",
}

// PitbullJob - represents a single Pitbull job.
type PitbullJob struct {
	common.BaseModel `bson:",inline"`

	InstanceId primitive.ObjectID `bson:"instanceId" json:"instanceId"`

	AcknowledgedAt  time.Time `bson:"acknowledgedAt" json:"acknowledgedAt"`
	LastScheduledAt time.Time `bson:"lastScheduledAt" json:"lastScheduledAt"`
	RejectedAt      time.Time `bson:"rejectedAt" json:"rejectedAt"`

	Status          JobStatus `bson:"status" json:"status"`
	RescheduleCount int       `bson:"rescheduleCount" json:"rescheduleCount"`
}

// NewPitbullJob - returns new PitbullJob instance.
func NewPitbullJob(instanceId primitive.ObjectID) *PitbullJob {
	job := &PitbullJob{
		InstanceId: instanceId,
		Status:     Scheduled,
	}

	job.ID = primitive.NewObjectID()

	return job
}

// FirstScheduledAt - "scheduled" timestamp is basically the same as
// "created" timestamp.
func (pj *PitbullJob) FirstScheduledAt() time.Time {
	return pj.CreatedAt
}
