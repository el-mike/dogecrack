package crack

import (
	"os"
	"runtime/debug"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull"
)

// JobDispatcher - observes redis-based PitbullQueue with BLPOP
// and delegates Pitbull runs to worker threads.
type JobDispatcher struct {
	jobRunner  *JobRunner
	jobManager *JobManager
	jobQueue   *JobQueue

	pollInterval time.Duration

	done chan bool

	logger *common.Logger
}

// NewJobDispatcher - returns new JobDispatcher.
func NewJobDispatcher(instanceManager *pitbull.InstanceManager, jobRunner *JobRunner, pollInterval time.Duration) *JobDispatcher {
	return &JobDispatcher{
		jobRunner:  jobRunner,
		jobManager: NewJobManager(instanceManager),

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

	defer func() {
		if r := recover(); r != nil {
			rd.logger.Err.Printf("Recovering from panic. reason: %v\n", r)
			rd.logger.Err.Printf("Stack: \n%s\n", string(debug.Stack()))
		}

		go rd.Start()
	}()

	for {
		select {
		case <-ticker.C:
			reason, err := rd.jobManager.CanDispatch()
			if err != nil {
				rd.logger.Err.Printf("Checking if job can be dispatched failed. reason: %v\n", err)
				continue
			}

			if reason != "" {
				rd.logger.Info.Printf("Job dispatch skipped. Reason: %s\n", reason)
				continue
			}

			job, err := rd.jobManager.DequeueJob()
			if err != nil {
				rd.logger.Err.Printf("Dequeueing job failed. reason: %v\n", err)

				continue
			}

			if job != nil {
				rd.logger.Info.Printf("job '%s' dequeued. Starting Pitbull...\n", job.ID.Hex())

				rd.jobRunner.Run(job)
			}

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
