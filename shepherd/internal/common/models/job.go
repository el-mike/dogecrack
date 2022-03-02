package models

// CrackJobStatusEnum - describes the status of an arbitrary Job.
type JobStatusEnum int8

// JobStatus - "namespaced" struct containing JobStatus values.
var JobStatus = struct {
	Scheduled    JobStatusEnum
	Processing   JobStatusEnum
	Rescheduled  JobStatusEnum
	Rejected     JobStatusEnum
	Acknowledged JobStatusEnum
}{
	Scheduled:    0,
	Processing:   1,
	Rescheduled:  2,
	Rejected:     3,
	Acknowledged: 4,
}

var jobStatusByName = map[JobStatusEnum]string{
	JobStatus.Scheduled:    "SCHEDULED",
	JobStatus.Processing:   "PROCESSING",
	JobStatus.Rescheduled:  "RESCHEDULED",
	JobStatus.Rejected:     "REJECTED",
	JobStatus.Acknowledged: "ACKNOWLEDGED",
}

// Job - represents a single, arbitrary Job. Job should be extended by more specific types.
type Job struct {
	BaseModel `bson:",inline"`

	FirstScheduledAt NullableTime `bson:"firstScheduledAt" json:"firstScheduledAt,omitempty"`
	LastScheduledAt  NullableTime `bson:"lastScheduledAt" json:"lastScheduledAt,omitempty"`
	AcknowledgedAt   NullableTime `bson:"acknowledgedAt" json:"acknowledgedAt,omitempty"`
	RejectedAt       NullableTime `bson:"rejectedAt" json:"rejectedAt,omitempty"`

	Status          JobStatusEnum `bson:"status" json:"status"`
	RescheduleCount int           `bson:"rescheduleCount" json:"rescheduleCount"`
}
