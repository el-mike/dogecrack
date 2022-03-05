package models

// InstancesStatistics - a struct containing PitbullInstances statistics.
type PitbullInstancesStatistics struct {
	All            int `bson:"all" json:"all"`
	WaitingForHost int `bson:"waitingForHost" json:"waitingForHost"`
	HostStarting   int `bson:"hostStarting" json:"hostStarting"`
	Running        int `bson:"running" json:"running"`
	Completed      int `bson:"completed" json:"completed"`
	Failed         int `bson:"failed" json:"failed"`
	Interrupted    int `bson:"interrupted" json:"interrupted"`
	Success        int `bson:"success" json:"success"`

	PasswordsChecked int     `bson:"passwordsChecked" json:"passwordsChecked"`
	TotalCost        float64 `bson:"totalCost" json:"totalCost"`
	AverageCost      float64 `bson:"averageCost" json:"averageCost"`
}
