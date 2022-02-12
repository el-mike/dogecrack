package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JobStatus - describes the status of a single PitbullJob.
type JobStatus int8

const (
	JobScheduled JobStatus = iota
	JobProcessing
	JobRescheduled
	JobRejected
	JobAcknowledged
)

var nameByJobStatus = map[JobStatus]string{
	JobScheduled:    "JOB_SCHEDULED",
	JobProcessing:   "JOB_PROCESSING",
	JobRescheduled:  "JOB_RESCHEDULED",
	JobRejected:     "JOB_REJECTED",
	JobAcknowledged: "JOB_ACKNOWLEDGED",
}

// PitbullJob - represents a single Pitbull job.
type PitbullJob struct {
	BaseModel `bson:",inline"`

	WalletString string `bson:"walletString" json:"walletString"`
	PasslistUrl  string `bson:"passlistUrl" json:"passlistUrl"`

	InstanceId primitive.ObjectID `bson:"instanceId" json:"instanceId"`
	Instance   *PitbullInstance   `bson:"instance,omitempty" json:"instance"`

	StartedAt time.Time `bson:"startedAt" json:"startedAt"`

	FirstScheduledAt time.Time `bson:"firstScheduledAt" json:"firstScheduledAt"`
	LastScheduledAt  time.Time `bson:"lastScheduledAt" json:"lastScheduledAt"`
	AcknowledgedAt   time.Time `bson:"acknowledgedAt" json:"acknowledgedAt"`
	RejectedAt       time.Time `bson:"rejectedAt" json:"rejectedAt"`

	Status          JobStatus `bson:"status" json:"status"`
	RescheduleCount int       `bson:"rescheduleCount" json:"rescheduleCount"`
}

// NewPitbullJob - returns new PitbullJob instance.
func NewPitbullJob(passlistUrl, walletString string) *PitbullJob {
	job := &PitbullJob{
		WalletString: walletString,
		PasslistUrl:  passlistUrl,
		Status:       JobScheduled,
	}

	job.ID = primitive.NewObjectID()

	return job
}
