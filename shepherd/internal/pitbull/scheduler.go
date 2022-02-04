package pitbull

import (
	"os"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/models"
)

const INSTANCES_LIMIT = 5

// Scheduler - entity responsible for scheduling Pitbull runs via JobQueue.
type Scheduler struct {
	queue *JobQueue

	logger *common.Logger
}

// NewScheduler - returns new PitbullScheduler instance.
func NewScheduler() *Scheduler {
	return &Scheduler{
		queue: NewJobQueue(persist.GetRedisClient()),

		logger: common.NewLogger("Scheduler", os.Stdout, os.Stderr),
	}
}

// ScheduleRun - schedules a single Pitbull run. If instances limit is not reach yet,
// it will run it immediately.
func (sc *Scheduler) ScheduleRun(instance *models.PitbullInstance) error {
	job := models.NewPitbullJob(instance)

	sc.queue.Enqueue(job)

	sc.logger.Info.Printf("instance scheduled\n")

	return nil
}
