package crack

import (
	"os"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
)

// Scheduler - entity responsible for scheduling Pitbull runs via JobQueue.
type Scheduler struct {
	jobQueue *JobQueue

	logger *common.Logger
}

// NewScheduler - returns new PitbullScheduler instance.
func NewScheduler() *Scheduler {
	return &Scheduler{
		jobQueue: NewJobQueue(),

		logger: common.NewLogger("Scheduler", os.Stdout, os.Stderr),
	}
}

// ScheduleRun - schedules a single Pitbull run. If instances limit is not reach yet,
// it will run it immediately.
func (sc *Scheduler) ScheduleRun(job *models.CrackJob) error {
	if err := sc.jobQueue.Enqueue(job.ID.Hex()); err != nil {
		return err
	}

	sc.logger.Info.Printf("instance scheduled\n")

	return nil
}
