package pitbull

import (
	"os"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/models"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/repositories"
)

const INSTANCES_LIMIT = 5

// Scheduler - entity responsible for scheduling Pitbull runs via JobQueue.
type Scheduler struct {
	queue *JobQueue

	jobRepository *repositories.JobRepository

	logger *common.Logger
}

// NewScheduler - returns new PitbullScheduler instance.
func NewScheduler() *Scheduler {
	return &Scheduler{
		queue: NewJobQueue(persist.GetRedisClient()),

		jobRepository: repositories.NewJobRepository(),

		logger: common.NewLogger("Scheduler", os.Stdout, os.Stderr),
	}
}

// ScheduleRun - schedules a single Pitbull run. If instances limit is not reach yet,
// it will run it immediately.
func (sc *Scheduler) ScheduleRun(instance *models.PitbullInstance) (*models.PitbullJob, error) {
	job := models.NewPitbullJob(instance.ID)

	if err := sc.jobRepository.Create(job); err != nil {
		return nil, err
	}

	if err := sc.queue.Enqueue(job.ID.Hex()); err != nil {
		return nil, err
	}

	sc.logger.Info.Printf("instance scheduled\n")

	return job, nil
}
