package pitbull

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/el-mike/dogecrack/shepherd/models"
	"github.com/robfig/cron/v3"
)

const INSTANCES_LIMIT = 5
const START_HOST_ATTEMPTS_LIMIT = 10
const CHECK_STATUS_RETRY_LIMIT = 10

// PitbullScheduler - entity responsible for monitoring and maintanence jobs.
type PitbullScheduler struct {
	pitbullManager *PitbullManager

	instancesLimit int
	instancesQueue *models.StringQueue

	jobsPool *jobsPool
}

// NewPitbullScheduler - returns new PitbullScheduler instance.
func NewPitbullScheduler(pitbullManager *PitbullManager) *PitbullScheduler {
	return &PitbullScheduler{
		pitbullManager: pitbullManager,

		instancesLimit: INSTANCES_LIMIT,
		instancesQueue: models.NewStringQueue(),

		jobsPool: newJobsPool(),
	}
}

// ScheduleRun - schedules a single Pitbull run. If instances limit is not reach yet,
// it will run it immediately.
func (ps *PitbullScheduler) ScheduleRun(instance *models.PitbullInstance) error {
	jobName := "scheduleRun"

	// If the queue is empty, we want to wake the dequeueJob.
	if ps.instancesQueue.IsEmpty() {
		if err := ps.dequeueJob(); err != nil {
			return err
		}
	}

	instanceId := instance.ID.Hex()
	ps.instancesQueue.Enqueue(instanceId)

	infoLogger := ps.getInfoJobLogger(jobName, instanceId)

	infoLogger.Printf("instance scheduled\n")

	return nil
}

func (ps *PitbullScheduler) dequeueJob() error {
	jobName := "dequeueJob"

	infoLogger := ps.getInfoJobLogger(jobName, "")
	errorLogger := ps.getErrorJobLogger(jobName, "")

	infoLogger.Printf("starting\n")

	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("*/30 * * * * *", func() {
		if !ps.instancesQueue.IsEmpty() && ps.jobsPool.Size() < INSTANCES_LIMIT {
			instanceId, err := ps.instancesQueue.Dequeue()
			if err != nil {
				errorLogger.Printf("instance dequeue failed. Reason: %s\n", err)
			}

			infoLogger.Printf("'%s' instance dequeued\n", instanceId)

			ps.startHostJob(instanceId)
		}

		// If queue has been emptied, we want to stop dequeueJob, to not waste resources
		// and prevent memory leaks. It will be waken up by ScheduleRun if needed.
		if ps.instancesQueue.IsEmpty() {
			infoLogger.Printf("queue empty, stopping\n")

			c.Stop()
			return
		}
	})

	if err != nil {
		errorLogger.Printf("cronjob error. Reason: %s\n", err)
		return err
	}

	c.Start()

	return nil
}

func (ps *PitbullScheduler) startHostJob(instanceId string) error {
	jobName := "startHostJob"

	infoLogger := ps.getInfoJobLogger(jobName, instanceId)
	errorLogger := ps.getErrorJobLogger(jobName, instanceId)

	infoLogger.Printf("starting host\n")

	_, err := ps.pitbullManager.RunHostForInstance(instanceId)
	if err != nil {
		errorLogger.Printf("starting host failed. Reason: %s\n", err)

		return err
	}

	c := cron.New(cron.WithSeconds())

	attemptsCount := 0

	// Runs at every 30th second.
	_, err = c.AddFunc("*/30 * * * * *", func() {
		attemptsCount += 1

		if attemptsCount >= START_HOST_ATTEMPTS_LIMIT {
			infoLogger.Printf("attempts limit reached, stopping job and host\n")

			if err := ps.pitbullManager.StopHostInstance(instanceId); err != nil {
				errorLogger.Printf("stopping host instance failed. Reason: %s\n", err)
			}

			c.Stop()
			return
		}

		instance, err := ps.pitbullManager.SyncInstance(instanceId)
		// Double check - if for some reason SyncInstance returned nil error and nil instance,
		// we want to return, to prevent nil pointer dereference.
		if err != nil || instance == nil {
			if err == nil {
				err = errors.New("instance is nil")
			}

			errorLogger.Printf("instance sync failed. Reason: %s\n", err)

			return
		}

		infoLogger.Printf("host status: %s\n", instance.HostInstance.HostStatus().Formatted())

		if instance.Status == models.Waiting {
			infoLogger.Printf("host started\n")

			ps.runPitbullJob(instanceId)

			c.Stop()
			return
		}
	})

	if err != nil {
		errorLogger.Printf("cronjob error. Reason: %s\n", err)
		return err
	}

	c.Start()

	return nil
}

func (ps *PitbullScheduler) runPitbullJob(instanceId string) error {
	jobName := "runPitbullJob"

	infoLogger := ps.getInfoJobLogger(jobName, instanceId)
	errorLogger := ps.getErrorJobLogger(jobName, instanceId)

	infoLogger.Printf("starting Pitbull\n")

	if _, err := ps.pitbullManager.RunPitbull(instanceId); err != nil {
		errorLogger.Printf("starting Pitbull failed. Reason: %s\n", err)

		return err
	}

	// This cron does not use WithSecond(), therefore first place
	// is minutes, not seconds.
	c := cron.New()

	retryCount := 0

	_, err := c.AddFunc("* * * * *", func() {
		if retryCount >= CHECK_STATUS_RETRY_LIMIT {
			infoLogger.Printf("retries limit reached, stopping job and host\n")

			if err := ps.pitbullManager.StopHostInstance(instanceId); err != nil {
				errorLogger.Printf("stopping host instance failed. Reason: %s\n", err)
			}

			c.Stop()
			return
		}

		instance, err := ps.pitbullManager.SyncInstance(instanceId)
		if err != nil || instance == nil {
			if err == nil {
				err = errors.New("instance is nil")
			}

			errorLogger.Printf("Pitbull sync failed. Reason: %s\n", err)

			retryCount += 1

			return
		}

		// If we have reached this point, that means sync was succesful - therefore, we want to
		// reset retry counter.
		retryCount = 0

		infoLogger.Printf("[Process]: %s | %s\n", instance.Status.Formatted(), instance.Progress.Formatted())

		if instance.Status == models.Finished ||
			instance.Status == models.Success {
			infoLogger.Printf("pitbull finished, stopping host instance\n")

			output, err := ps.pitbullManager.GetInstanceOutput(instance)
			if err != nil {
				errorLogger.Printf("output retrieval failed. Reason: %s\n", err)
			}

			if err := ps.pitbullManager.StopHostInstance(instanceId); err != nil {
				errorLogger.Printf("stopping host instance '%d' failed. reason: %s\n", instance.HostInstance.ProviderId(), err)
			}

			infoLogger.Printf("host instance stopped \n")

			if output != "" {
				instance.LastOutput = output

				if err := ps.pitbullManager.UpdateInstance(instance); err != nil {
					errorLogger.Printf("saving last output failed. Reason: %s\n", err)
				}
			}

			c.Stop()

			infoLogger.Printf("job completed\n")
			return
		}
	})

	if err != nil {
		errorLogger.Printf("cronjob error. Reason: %s\n", err)
		return err
	}

	c.Start()

	return nil
}

// getJobLogger - returns a logger with predefined "tags" and correct setup.
func (ps *PitbullScheduler) getJobLogger(out io.Writer, logType string, jobName, instanceId string) *log.Logger {
	instanceTag := ""

	if instanceId != "" {
		instanceTag = fmt.Sprintf("[%s]", instanceId)
	}
	return log.New(out, fmt.Sprintf("[Scheduler][%s][%s]%s: ", logType, jobName, instanceTag), log.Ldate|log.Ltime|log.Lmsgprefix)
}

// getInfoJobLogger - returns info job logger.
func (ps *PitbullScheduler) getInfoJobLogger(jobName, instanceId string) *log.Logger {
	// TODO: Add file handling for output.
	return ps.getJobLogger(os.Stdout, "Info", jobName, instanceId)
}

// getErrorJobLogger - returns error job logger.
func (ps *PitbullScheduler) getErrorJobLogger(jobName, instanceId string) *log.Logger {
	// TODO: Add file handling for output.
	return ps.getJobLogger(os.Stderr, "Error", jobName, instanceId)
}

type jobsPool struct {
	sync.RWMutex

	instancesIds []string
}

func newJobsPool() *jobsPool {
	return &jobsPool{
		instancesIds: []string{},
	}
}

func (jp *jobsPool) AddJob(instanceId string) {
	jp.Lock()
	defer jp.Unlock()

	if jp.instancesIds == nil {
		jp.instancesIds = []string{}
	}

	jp.instancesIds = append(jp.instancesIds, instanceId)
}

func (jp *jobsPool) RemoveJob(instanceId string) {
	jp.Lock()
	defer jp.Unlock()

	if len(jp.instancesIds) == 0 {
		return
	}

	result := []string{}

	for _, id := range jp.instancesIds {
		if id != instanceId {
			result = append(result, id)
		}
	}

	jp.instancesIds = result
}

func (jp *jobsPool) Size() int {
	return len(jp.instancesIds)
}
