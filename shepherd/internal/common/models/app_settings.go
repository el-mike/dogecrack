package models

// AppSettings - contains all app's specific settings, that should be editable by admin.
type AppSettings struct {
	BaseModel `bson:",inline"`

	StartHostAttemptsLimit int8 `bson:"startHostAttemptsLimit" json:"startHostAttemptsLimit"`
	CheckStatusRetryLimit  int8 `bson:"checkStatusRetryLimit" json:"checkStatusRetryLimit"`
	StalledProgressLimit   int8 `bson:"stalledProgressLimit" json:"stalledProgressLimit"`
	RescheduleLimit        int8 `bson:"rescheduleLimit" json:"rescheduleLimit"`
	RunningInstancesLimit  int8 `bson:"runningInstancesLimit" json:"runningInstancesLimit"`

	CheckHostInterval    int16 `bson:"checkHostInterval" json:"checkHostInterval"`
	CheckPitbullInterval int16 `bson:"checkPitbullInterval" json:"checkPitbullInterval"`
}

// NewAppSettings - returns new AppSettings instance.
func NewAppSettings() *AppSettings {
	return &AppSettings{}
}

// Merge - merges source AppSettings into current receiver.
func (as *AppSettings) Merge(source *AppSettings) {
	if source.StartHostAttemptsLimit != 0 {
		as.StartHostAttemptsLimit = source.StartHostAttemptsLimit
	}

	if source.CheckStatusRetryLimit != 0 {
		as.CheckStatusRetryLimit = source.CheckStatusRetryLimit
	}

	if source.StalledProgressLimit != 0 {
		as.StalledProgressLimit = source.StalledProgressLimit
	}

	if source.RescheduleLimit != 0 {
		as.RescheduleLimit = source.RescheduleLimit
	}

	if source.RunningInstancesLimit != 0 {
		as.RunningInstancesLimit = source.RunningInstancesLimit
	}

	if source.CheckHostInterval != 0 {
		as.CheckHostInterval = source.CheckHostInterval
	}

	if source.CheckPitbullInterval != 0 {
		as.CheckPitbullInterval = source.CheckPitbullInterval
	}
}

// AppSettingsPayload - an alias type for possible AppSettings payload.
type AppSettingsPayload = AppSettings
