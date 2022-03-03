package models

// JobsStatistics - a struct containing counts of jobs according to statuses.
type JobsStatistics struct {
	All          int `bson:"all" json:"all"`
	Acknowledged int `bson:"acknowledged" json:"acknowledged"`
	Processing   int `bson:"processing" json:"processing"`
	Queued       int `bson:"queued" json:"queued"`
	Rejected     int `bson:"rejected" json:"rejected"`
}
