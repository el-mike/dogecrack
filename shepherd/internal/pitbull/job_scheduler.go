package pitbull

import (
	"os"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
)

const INSTANCES_LIMIT = 5

// Scheduler - entity responsible for scheduling Pitbull runs via JobQueue.
type Scheduler struct {
	jobQueue *JobQueue

	jobManager *JobManager

	logger *common.Logger
}

// NewScheduler - returns new PitbullScheduler instance.
func NewScheduler(instanceManager *InstanceManager) *Scheduler {
	return &Scheduler{
		jobQueue:   NewJobQueue(),
		jobManager: NewJobManager(instanceManager),

		logger: common.NewLogger("Scheduler", os.Stdout, os.Stderr),
	}
}

// ScheduleRun - schedules a single Pitbull run. If instances limit is not reach yet,
// it will run it immediately.
func (sc *Scheduler) ScheduleRun(passlistUrl, walletString string) (*models.PitbullJob, error) {
	job, err := sc.jobManager.CreateJob(passlistUrl, walletString)
	if err != nil {
		return nil, err
	}

	if err := sc.jobQueue.Enqueue(job.ID.Hex()); err != nil {
		return nil, err
	}

	sc.logger.Info.Printf("instance scheduled\n")

	return job, nil
}
