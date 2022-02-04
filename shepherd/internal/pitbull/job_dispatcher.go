package pitbull

import (
	"os"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"github.com/el-mike/dogecrack/shepherd/internal/utils"
)

// JobDispatcher - observes redis-based PitbullQueue with BLPOP
// and delegates Pitbull runs to worker threads.
type JobDispatcher struct {
	runner *Runner

	queue *JobQueue
	done  chan bool

	pollInterval time.Duration

	logger *utils.Logger
}

// NewJobDispatcher - returns new JobDispatcher.
func NewJobDispatcher(runner *Runner, pollInterval time.Duration) *JobDispatcher {
	return &JobDispatcher{
		runner: runner,

		queue: NewJobQueue(persist.GetRedisClient()),
		done:  make(chan bool),

		pollInterval: pollInterval,

		logger: utils.NewLogger("Dispatcher", os.Stdout, os.Stderr),
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

			job, err := rd.queue.Dequeue()
			if err != nil {
				rd.logger.Err.Println(err)

				continue
			}

			if job == nil {
				rd.logger.Info.Println("No scheduled jobs found in the queue.")

				continue
			}

			rd.logger.Info.Printf("job for instance '%s' dequeued. Starting Pitbull...\n", job.InstanceId)

			rd.runner.Run(job.InstanceId)

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
