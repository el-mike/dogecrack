package pitbull

import (
	"os"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/models"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/repositories"
)

// JobDispatcher - observes redis-based PitbullQueue with BLPOP
// and delegates Pitbull runs to worker threads.
type JobDispatcher struct {
	runner        *Runner
	jobRepository *repositories.JobRepository

	queue *JobQueue
	done  chan bool

	pollInterval time.Duration

	logger *common.Logger
}

// NewJobDispatcher - returns new JobDispatcher.
func NewJobDispatcher(runner *Runner, pollInterval time.Duration) *JobDispatcher {
	return &JobDispatcher{
		runner:        runner,
		jobRepository: repositories.NewJobRepository(),

		queue: NewJobQueue(persist.GetRedisClient()),
		done:  make(chan bool),

		pollInterval: pollInterval,

		logger: common.NewLogger("Dispatcher", os.Stdout, os.Stderr),
	}
}

// Start - starts JobDispatcher.
func (rd *JobDispatcher) Start() {
	rd.logger.Info.Printf("started. Poll interval: %v\n", rd.pollInterval)

	ticker := time.NewTicker(rd.pollInterval)

	for {
		select {
		case <-ticker.C:
			rd.logger.Info.Println("Checking for scheduled jobs...")

			jobId, err := rd.queue.Dequeue()
			if err != nil {
				rd.logger.Err.Println(err)

				continue
			}

			if jobId == "" {
				rd.logger.Info.Println("No scheduled jobs found in the queue.")

				continue
			}

			job, err := rd.jobRepository.GetById(jobId)
			if err != nil {
				rd.logger.Err.Println(err)

				continue
			}

			job.Status = models.Processing
			job.LastScheduledAt = time.Now()

			if err := rd.jobRepository.Update(job); err != nil {
				rd.logger.Err.Printf("updating status for job '%s' failed. reason: %s\n", job.ID.Hex(), err)
			}

			rd.logger.Info.Printf("job '%s' dequeued. Starting Pitbull...\n", job.ID.Hex())

			rd.runner.Run(job)

		case <-rd.done:
			ticker.Stop()
			rd.logger.Info.Println("stopped.")

			return
		}
	}
}

// Stop - stops JobDispatcher.
func (rd *JobDispatcher) Stop() {
	rd.done <- true
}
