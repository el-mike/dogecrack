package models

// AppStatistics - represents a general, application statistics.
type AppStatistics struct {
	CrackJobs        *CrackJobsStatistics        `json:"crackJobs"`
	PitbullInstances *PitbullInstancesStatistics `json:"pitbullInstances"`
}
