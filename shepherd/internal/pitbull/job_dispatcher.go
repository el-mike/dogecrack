package pitbull

import (
	"os"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/models"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/repositories"
)

// JobDispatcher - observes redis-based PitbullQueue with BLPOP
// and delegates Pitbull runs to worker threads.
type JobDispatcher struct {
	jobRunner     *JobRunner
	jobRepository *repositories.JobRepository
	jobQueue      *JobQueue

	pollInterval time.Duration

	done chan bool

	logger *common.Logger
}

// NewJobDispatcher - returns new JobDispatcher.
func NewJobDispatcher(jobRunner *JobRunner, pollInterval time.Duration) *JobDispatcher {
	return &JobDispatcher{
		jobRunner:     jobRunner,
		jobRepository: repositories.NewJobRepository(),

		pollInterval: pollInterval,

		jobQueue: NewJobQueue(),
		done:     make(chan bool),

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

			jobId, err := rd.jobQueue.Dequeue()
			if err != nil {
				rd.logger.Err.Println(err)

				continue
			}

			if jobId == "" {
				continue
			}

			job, err := rd.jobRepository.GetById(jobId)
			if err != nil {
				rd.logger.Err.Println(err)

				continue
			}

			job.Status = models.Processing
			job.StartedAt = time.Now()

			if err := rd.jobRepository.Update(job); err != nil {
				rd.logger.Err.Printf("updating status for job '%s' failed. reason: %s\n", job.ID.Hex(), err)
			}

			rd.logger.Info.Printf("job '%s' dequeued. Starting Pitbull...\n", job.ID.Hex())

			rd.jobRunner.Run(job)

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
