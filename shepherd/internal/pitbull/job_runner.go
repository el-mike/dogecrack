package pitbull

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	startHostAttemptsLimit = 10
	checkStatusRetryLimit  = 10
	stalledProgressLimit   = 10
	rescheduleLimit        = 5
	checkHostInterval      = 15 * time.Second
	checkPitbullInterval   = 30 * time.Second
)

// JobRunner - entity responsible for running and monitoring Pitbull jobs.
type JobRunner struct {
	instanceManager *InstanceManager
	jobQueue        *JobQueue
	jobManager      *JobManager
}

// NewJobRunner - returns new PitbullRunner instance.
func NewJobRunner(instanceManager *InstanceManager) *JobRunner {
	return &JobRunner{
		instanceManager: instanceManager,
		jobQueue:        NewJobQueue(),
		jobManager:      NewJobManager(instanceManager),
	}
}

// Run - starts single Pitbull run.
func (ru *JobRunner) Run(job *models.PitbullJob) {
	go ru.runSingle(job)
}

func (ru *JobRunner) runSingle(job *models.PitbullJob) {
	defer func() {
		if r := recover(); r != nil {
			logger := common.NewLogger("Runner", os.Stdout, os.Stderr, "recovery", job.ID.Hex())

			logger.Err.Printf("Recovering from panic. reason: %v\n", r)
			logger.Err.Printf("Stack: \n%s\n", string(debug.Stack()))
		}
	}()

	ru.assignInstance(job)
}

// createInstance - creates and attaches a new PitbullInstance to PitbullJob.
// If rescheduled, it tries to destroy the previous HostInstance.
func (ru *JobRunner) assignInstance(job *models.PitbullJob) {
	logger := common.NewLogger("Runner", os.Stdout, os.Stderr, "assignInstance", job.ID.Hex())

	if job.InstanceId != primitive.NilObjectID {
		instance, err := ru.instanceManager.GetInstanceById(job.InstanceId.Hex())
		if err != nil {
			logger.Err.Printf("Retrieving previous instance failed. reason: %v\n", err)
		} else {
			if instance != nil && instance.Active() {
				logger.Info.Printf("Stopping previously assigned instance...\n")

				if err := ru.instanceManager.StopHostInstance(job.InstanceId.Hex()); err != nil {
					logger.Err.Printf("Stopping previously assigned instance failed. reason: %v\n", err)
				}
			}
		}
	}

	job, err := ru.jobManager.AssignInstance(job)
	if err != nil {
		logger.Err.Printf("Assign instance failed. reason: %v\n", err)

		ru.handleFailure(job, logger.DecorateErr(err))
		return
	}

	ru.startHost(job)
}

// startHost - starts a single host for Pitbull process to work in.
func (ru *JobRunner) startHost(job *models.PitbullJob) {
	logger := common.NewLogger("Runner", os.Stdout, os.Stderr, "startHost", job.ID.Hex())

	logger.Info.Println("starting host.")

	_, err := ru.instanceManager.RunHostForInstance(job.InstanceId.Hex())
	if err != nil {
		logger.Err.Printf("starting host failed. Reason: %s\n", err)

		ru.handleFailure(job, logger.DecorateErr(err))
		return
	}

	ticker := time.NewTicker(checkHostInterval)

	attemptsCount := 0

	for range ticker.C {
		attemptsCount += 1

		if attemptsCount >= startHostAttemptsLimit {
			logger.Err.Printf("attempts limit reached, stopping job and host\n")

			if err := ru.instanceManager.StopHostInstance(job.InstanceId.Hex()); err != nil {
				logger.Err.Printf("stopping host instance failed. Reason: %s\n", err)
			}

			ticker.Stop()

			ru.handleFailure(job, logger.DecorateErr(err))
			return
		}

		instance, err := ru.instanceManager.SyncInstance(job.InstanceId.Hex())
		// Double check - if for some reason SyncInstance returned nil error and nil instance,
		// we want to return, to prevent nil pointer dereference.
		if err != nil || instance == nil {
			if err == nil {
				err = errors.New("instance is nil")
			}

			logger.Err.Printf("instance sync failed. Reason: %s\n", err)

			continue
		}

		logger.Info.Printf("host status: %s\n", instance.HostInstance.HostStatus().Formatted())

		if instance.Status == models.Waiting {
			logger.Info.Printf("host started\n")

			go ru.runPitbull(job)

			ticker.Stop()
			break
		}
	}
}

func (ru *JobRunner) runPitbull(job *models.PitbullJob) {
	logger := common.NewLogger("Runner", os.Stdout, os.Stderr, "runPitbull", job.ID.Hex())

	logger.Info.Printf("starting Pitbull\n")

	if _, err := ru.instanceManager.RunPitbull(job.InstanceId.Hex()); err != nil {
		logger.Err.Printf("starting Pitbull failed. Reason: %s\n", err)

		ru.handleFailure(job, logger.DecorateErr(err))
		return
	}

	ticker := time.NewTicker(checkPitbullInterval)

	retryCount := 0

	lastProgress := int64(0)
	stalledProgressCount := 0

	for range ticker.C {
		if retryCount >= checkStatusRetryLimit {
			err := fmt.Errorf("retries limit reached")

			logger.Err.Printf("%s, stopping job and host\n", err)

			if err := ru.instanceManager.StopHostInstance(job.InstanceId.Hex()); err != nil {
				logger.Err.Printf("stopping host instance failed. Reason: %s\n", err)
			}

			ticker.Stop()

			ru.handleFailure(job, logger.DecorateErr(err))
			return
		}

		instance, err := ru.instanceManager.SyncInstance(job.InstanceId.Hex())
		if err != nil || instance == nil {
			if err == nil {
				err = errors.New("instance is nil")
			}

			logger.Err.Printf("Pitbull sync failed. Reason: %s\n", err)

			retryCount += 1

			continue
		}

		// If we have reached this point, that means sync was succesful - therefore, we want to
		// reset retry counter.
		retryCount = 0

		logger.Info.Printf("[Process]: %s | %s\n", instance.Status.Formatted(), instance.Progress.Formatted())

		if instance.Status == models.Finished ||
			instance.Status == models.Success {
			logger.Info.Printf("pitbull finished, stopping host instance\n")

			output, err := ru.instanceManager.GetInstanceOutput(instance)
			if err != nil {
				logger.Err.Printf("output retrieval failed. Reason: %s\n", err)
			}

			if err := ru.instanceManager.StopHostInstance(job.InstanceId.Hex()); err != nil {
				logger.Err.Printf("stopping host instance '%d' failed. reason: %s\n", instance.HostInstance.ProviderId(), err)
			}

			logger.Info.Printf("host instance stopped \n")

			if output != "" {
				instance.LastOutput = output

				if err := ru.instanceManager.UpdateInstance(instance); err != nil {
					logger.Err.Printf("saving last output failed. Reason: %s\n", err)
				}
			}

			logger.Info.Printf("job completed\n")

			ticker.Stop()
			ru.handleCompletion(job)

			return
		}

		currentProgress := instance.Progress.Checked

		// If progress did not change since the last iteration, we increment
		// the counter. Otherwise we want to reset it, since progress has been made.
		if currentProgress == lastProgress {
			stalledProgressCount += 1
		} else {
			stalledProgressCount = 0
		}

		if stalledProgressCount >= stalledProgressLimit {
			logger.Err.Printf("pitbull progress stalled, stopping job and host\n")

			if err := ru.instanceManager.StopHostInstance(job.InstanceId.Hex()); err != nil {
				logger.Err.Printf("stopping host instance failed. Reason: %s\n", err)
			}

			ticker.Stop()

			ru.handleFailure(job, logger.DecorateErr(err))
			return
		}

		lastProgress = instance.Progress.Checked

	}
}

// handleCompletion - performs any cleanups and updates after completing the job.
func (ru *JobRunner) handleCompletion(job *models.PitbullJob) {
	logger := common.NewLogger("Runner", os.Stdout, os.Stderr, "cleanup", job.ID.Hex())

	if err := ru.jobManager.AcknowledgeJob(job); err != nil {
		logger.Err.Printf("Acknowledge failed. reason: %s\n", err)

		return
	}

	logger.Info.Printf("Job acknowledged\n")
}

// handleFailure - handles a failure scenario by rescheduling or rejecting the job,
// based on its history and status.
func (ru *JobRunner) handleFailure(job *models.PitbullJob, reason error) {
	logger := common.NewLogger("Runner", os.Stdout, os.Stderr, "cleanup", job.ID.Hex())

	logger.Info.Printf("Failure handling started\n")

	if job.RescheduleCount > rescheduleLimit {
		logger.Info.Printf("Reschedule limit reached, rejecting\n")

		if err := ru.jobManager.RejectJob(job, reason); err != nil {
			logger.Err.Printf("Rejecting failed. reason: %s\n", err)

			return
		}
	} else {
		logger.Info.Printf("Rescheduling.\n")

		if err := ru.jobManager.RescheduleJob(job, reason); err != nil {
			logger.Err.Printf("Rescheduling failed. reason: %s\n", err)

			return
		}
	}
}
