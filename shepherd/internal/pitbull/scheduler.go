package pitbull

import (
	"os"

	"github.com/el-mike/dogecrack/shepherd/internal/models"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"github.com/el-mike/dogecrack/shepherd/internal/utils"
)

const INSTANCES_LIMIT = 5

// PitbullScheduler - entity responsible for scheduling Pitbull runs via JobQueue.
type PitbullScheduler struct {
	pitbullManager *PitbullManager

	queue *JobQueue

	logger *utils.Logger
}

// NewPitbullScheduler - returns new PitbullScheduler instance.
func NewPitbullScheduler(pitbullManager *PitbullManager) *PitbullScheduler {
	return &PitbullScheduler{
		pitbullManager: pitbullManager,

		queue: NewJobQueue(persist.GetRedisClient()),

		logger: utils.NewLogger("Scheduler", os.Stdout, os.Stderr),
	}
}

// ScheduleRun - schedules a single Pitbull run. If instances limit is not reach yet,
// it will run it immediately.
func (ps *PitbullScheduler) ScheduleRun(instance *models.PitbullInstance) error {
	job := models.NewPitbullJob(instance)

	ps.queue.Enqueue(job)

	ps.logger.Info.Printf("instance scheduled\n")

	return nil
}
