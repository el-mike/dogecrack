package repositories

import (
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"go.mongodb.org/mongo-driver/mongo"
)

const jobsCollection = "jobs"

// JobRepository - MongoDB-backed repository for handling Pitbull jobs.
type JobRepository struct {
	db *mongo.Database
}

// NewJobRepository - returns new PitbullJobRepository instance.
func NewJobRepository() *JobRepository {
	return &JobRepository{
		db: persist.GetDatabase(),
	}
}
