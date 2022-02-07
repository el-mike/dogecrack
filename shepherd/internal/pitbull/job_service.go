package pitbull

import (
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/models"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/repositories"
)

// JobService - simple facade for operations on PitbullJobs.
type JobService struct {
	jobRepository *repositories.JobRepository
}

// NewJobService - returns new JobService instance.
func NewJobService() *JobService {
	return &JobService{
		jobRepository: repositories.NewJobRepository(),
	}
}

func (js *JobService) GetJobs() ([]*models.PitbullJob, error) {
	return nil, nil
}
