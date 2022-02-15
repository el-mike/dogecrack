package models

// AppEnums - contains all app's enums in a form of map of name->value maps.
type AppEnums struct {
	JobStatus             map[string]JobStatus     `json:"jobStatus"`
	PitbullInstanceStatus map[string]PitbullStatus `json:"pitbullInstanceStatus"`
}

// GetAppEnums - builds and returns an instance of AppEnums.
func GetAppEnums() *AppEnums {
	return &AppEnums{
		JobStatus: map[string]JobStatus{
			"JOB_SCHEDULED":    JobScheduled,
			"JOB_PROCESSING":   JobProcessing,
			"JOB_RESCHEDULED":  JobRescheduled,
			"JOB_REJECTED":     JobRejected,
			"JOB_ACKNOWLEDGED": JobAcknowledged,
		},
		PitbullInstanceStatus: map[string]PitbullStatus{
			"WAITING_FOR_HOST": WaitingForHost,
			"HOST_STARTING":    HostStarting,
			"WAITING":          Waiting,
			"RUNNING":          Running,
			"FINISHED":         Finished,
			"SUCCESS":          Success,
			"INTERRUPTED":      Interrupted,
		},
	}
}
