package pitbull

import (
	"os"
	"runtime/debug"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
)

// InstanceCollector - an entity responsible for detecting and disposing
// orphaned Pitbull instances (the ones with inactive parent PitbullJob).
type InstanceCollector struct {
	instanceManager *InstanceManager

	interval time.Duration

	done chan bool

	logger *common.Logger
}

// NewInstanceCollector - returns new Collector instance.
func NewInstanceCollector(instanceManager *InstanceManager, interval time.Duration) *InstanceCollector {
	return &InstanceCollector{
		instanceManager: instanceManager,

		interval: interval,

		done: make(chan bool),

		logger: common.NewLogger("Collector", os.Stdout, os.Stderr),
	}
}

// Start - starts Collector.
func (cl *InstanceCollector) Start() {
	cl.logger.Info.Printf("started. Interval: %v\n", cl.interval)

	defer func() {
		if r := recover(); r != nil {
			cl.logger.Err.Printf("Recovering from panic. reason: %v\n", r)
			cl.logger.Err.Printf("Stack: \n%s\n", string(debug.Stack()))
		}

		go cl.Start()
	}()

	ticker := time.NewTicker(cl.interval)

	for {
		select {
		case <-ticker.C:
			cl.logger.Info.Println("Checking for orphan instances...")

			instances, err := cl.instanceManager.GetOrphanInstances()
			if err != nil {
				cl.logger.Err.Printf("Getting orphan instances failed. reason: %v\n", err)

				continue
			}

			if instances == nil || len(instances) == 0 {
				continue
			}

			for _, instance := range instances {
				if instance == nil {
					continue
				}

				if err := cl.instanceManager.StopHostInstance(instance.ID.Hex()); err != nil {
					cl.logger.Err.Printf("Stopping orphan instance failed. reason: %v\n", err)

					continue
				}

				cl.logger.Info.Printf("Instance '%s' has been stopped.\n", instance.ID.Hex())

				if err := cl.instanceManager.MarkInstanceAsInterrupted(instance); err != nil {
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
func (cl *InstanceCollector) Stop() {
	cl.done <- true
}
