package crack

import (
	"os"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull"
)

const INSTANCES_LIMIT = 5

// Scheduler - entity responsible for scheduling Pitbull runs via JobQueue.
type Scheduler struct {
	jobQueue *JobQueue

	jobManager *JobManager

	logger *common.Logger
}

// NewScheduler - returns new PitbullScheduler instance.
func NewScheduler(instanceManager *pitbull.InstanceManager) *Scheduler {
	return &Scheduler{
		jobQueue:   NewJobQueue(),
		jobManager: NewJobManager(instanceManager),

		logger: common.NewLogger("Scheduler", os.Stdout, os.Stderr),
	}
}

// ScheduleRun - schedules a single Pitbull run. If instances limit is not reach yet,
// it will run it immediately.
func (sc *Scheduler) ScheduleRun(keyword, passlistUrl, walletString string) (*models.CrackJob, error) {
	job, err := sc.jobManager.CreateJob(keyword, passlistUrl, walletString)
	if err != nil {
		return nil, err
	}

	if err := sc.jobQueue.Enqueue(job.ID.Hex()); err != nil {
		return nil, err
	}

	sc.logger.Info.Printf("instance scheduled\n")

	return job, nil
}
