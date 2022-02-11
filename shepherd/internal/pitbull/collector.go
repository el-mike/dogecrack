package pitbull

import (
	"os"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
)

// Collector - an entity responsible for detecting and disposing
// orphaned Pitbull instances (the ones with inactive parent PitbullJob).
type Collector struct {
	jobManager *JobManager
	queue      *JobQueue
	manager    *Manager

	interval time.Duration

	done chan bool

	logger *common.Logger
}

// NewCollector - returns new Collector instance.
func NewCollector(manager *Manager, interval time.Duration) *Collector {
	return &Collector{
		jobManager: NewJobManager(),
		queue:      NewJobQueue(),
		manager:    manager,

		interval: interval,

		done: make(chan bool),

		logger: common.NewLogger("Collector", os.Stdout, os.Stderr),
	}
}

// Start - starts Collector.
func (cl *Collector) Start() {
	cl.logger.Info.Printf("started. Interval: %v\n", cl.interval)

	defer func() {
		if r := recover(); r != nil {
			cl.logger.Err.Printf("Recovering from panic. reason: %v\n", r)
		}

		cl.Start()
	}()

	ticker := time.NewTicker(cl.interval)

	for {
		select {
		case <-ticker.C:
			cl.logger.Info.Println("Checking for orphan instances...")

			jobs, err := cl.jobManager.GetCompletedJobWithActiveInstance()
			if err != nil {
				cl.logger.Err.Println(err)

				continue
			}

			if jobs == nil || len(jobs) == 0 {
				continue
			}

			for _, job := range jobs {
				instance := job.Instance

				if instance == nil {
					continue
				}

				if err := cl.manager.StopHostInstance(instance.ID.Hex()); err != nil {
					cl.logger.Err.Printf("Stopping orphan instance failed. reason: %v\n", err)

					continue
				}

				cl.logger.Info.Printf("Instance '%s' has been stopped.\n", instance.ID.Hex())

				if err := cl.jobManager.MarkInstanceAsInterrupted(instance); err != nil {
					cl.logger.Err.Printf("Marking instance '%s' as interrupted failed. reason: %v\n", instance.ID.Hex(), err)
				}
			}

		case <-cl.done:
			ticker.Stop()

			cl.logger.Info.Println("stopped.")

			return
		}
	}
}

// Stop - stops Collector.
func (cl *Collector) Stop() {
	cl.done <- true
}
