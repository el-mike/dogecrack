package models

// AppEnums - contains all app's enums in a form of map of name->value maps.
type AppEnums struct {
	JobStatus             map[string]JobStatusEnum             `json:"jobStatus"`
	PitbullInstanceStatus map[string]PitbullInstanceStatusEnum `json:"pitbullInstanceStatus"`
	PitbullStatus         map[string]PitbullStatusEnum         `json:"pitbullStatus"`
}

// GetAppEnums - builds and returns an instance of AppEnums.
func GetAppEnums() *AppEnums {
	return &AppEnums{
		JobStatus: map[string]JobStatusEnum{
			"SCHEDULED":    JobStatus.Scheduled,
			"PROCESSING":   JobStatus.Processing,
			"RESCHEDULED":  JobStatus.Rescheduled,
			"REJECTED":     JobStatus.Rejected,
			"ACKNOWLEDGED": JobStatus.Acknowledged,
		},
		PitbullInstanceStatus: map[string]PitbullInstanceStatusEnum{
			"WAITING_FOR_HOST": PitbullInstanceStatus.WaitingForHost,
			"HOST_STARTING":    PitbullInstanceStatus.HostStarting,
			"RUNNING":          PitbullInstanceStatus.Running,
			"COMPLETED":        PitbullInstanceStatus.Completed,
			"INTERRUPTED":      PitbullInstanceStatus.Interrupted,
			"FAILED":           PitbullInstanceStatus.Failed,
			"SUCCESS":          PitbullInstanceStatus.Success,
		},
		PitbullStatus: map[string]PitbullStatusEnum{
			"WAITING":  PitbullStatus.Waiting,
			"RUNNING":  PitbullStatus.Running,
			"FINISHED": PitbullStatus.Finished,
			"SUCCESS":  PitbullStatus.Success,
		},
	}
}
