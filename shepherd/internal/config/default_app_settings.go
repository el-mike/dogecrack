package config

import "github.com/el-mike/dogecrack/shepherd/internal/common/models"

// DefaultAppSettings - default values for AppSettings.
var DefaultAppSettings = &models.AppSettings{
	StartHostAttemptsLimit: 10,
	CheckStatusRetryLimit:  10,
	StalledProgressLimit:   10,
	RescheduleLimit:        5,

	RunningInstancesLimit: 10,

	// Below values are seconds.
	CheckHostInterval:    15,
	CheckPitbullInterval: 30,

	// Seems to yield best result cost-wise.
	VastInstanceSearchCriteria: "num_gpus=4 dph<2.5",

	MinPasswordLength: 0,
	MaxPasswordLength: 18,
}
