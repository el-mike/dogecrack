package pitbull

import (
	"os"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/models"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/repositories"
)

// Collector - an entity responsible for detecting and disposing
// orphaned Pitbull instances (the ones with inactive parent PitbullJob).
type Collector struct {
	jobRepository      *repositories.JobRepository
	instanceRepository *repositories.InstanceRepository

	interval time.Duration

	done chan bool

	logger *common.Logger
}

// NewCollector - returns new Collector instance.
func NewCollector(interval time.Duration) *Collector {
	return &Collector{
		jobRepository:      repositories.NewJobRepository(),
		instanceRepository: repositories.NewInstanceRepository(),

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

			instances, err := cl.instanceRepository.GetActiveInstances()
			if err != nil {
				cl.logger.Err.Println(err)

				continue
			}

			if instances == nil || len(instances) == 0 {
				continue
			}

			jobs, err := cl.jobRepository.GetCompletedWithActiveInstance()
			if err != nil {
				cl.logger.Err.Println(err)

				continue
			}

			for _, job := range jobs {

				if job.Status != models.Processing {

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
