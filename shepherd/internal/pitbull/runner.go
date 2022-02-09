package pitbull

import (
	"errors"
	"os"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/models"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/repositories"
)

const (
	startHostAttemptsLimit = 10
	checkStatusRetryLimit  = 10
	stalledProgressLimit   = 10
	rescheduleLimit        = 5
	checkHostInterval      = 15 * time.Second
	checkPitbullInterval   = 30 * time.Second
)

// Runner - entity responsible for running and monitoring Pitbull jobs.
type Runner struct {
	manager       *Manager
	queue         *JobQueue
	jobRepository *repositories.JobRepository
}

// NewRunner - returns new PitbullRunner instance.
func NewRunner(manager *Manager) *Runner {
	return &Runner{
		manager:       manager,
		queue:         NewJobQueue(persist.GetRedisClient()),
		jobRepository: repositories.NewJobRepository(),
	}
}

// Run - starts single Pitbull run.
func (ru *Runner) Run(job *models.PitbullJob) {
	go ru.runSingle(job)
}

func (ru *Runner) runSingle(job *models.PitbullJob) {
	defer func() {
		if r := recover(); r != nil {
			logger := common.NewLogger("Runner", os.Stdout, os.Stderr, "recovery", job.ID.Hex())
			logger.Err.Printf("Recovering from panic. reason: %v\n", r)
		}
	}()

	ru.startHost(job)
}

// startHost - starts a single host for Pitbull process to work in.
func (ru *Runner) startHost(job *models.PitbullJob) {
	logger := common.NewLogger("Runner", os.Stdout, os.Stderr, "startHost", job.ID.Hex())

	logger.Info.Println("starting host.")

	_, err := ru.manager.RunHostForInstance(job.InstanceId.Hex())
	if err != nil {
		logger.Err.Printf("starting host failed. Reason: %s\n", err)

		ru.handleFailure(job)
		return
	}

	ticker := time.NewTicker(checkHostInterval)

	attemptsCount := 0

	for range ticker.C {
		attemptsCount += 1

		if attemptsCount >= startHostAttemptsLimit {
			logger.Err.Printf("attempts limit reached, stopping job and host\n")

			if err := ru.manager.StopHostInstance(job.InstanceId.Hex()); err != nil {
				logger.Err.Printf("stopping host instance failed. Reason: %s\n", err)
			}

			ticker.Stop()

			ru.handleFailure(job)
			return
		}

		instance, err := ru.manager.SyncInstance(job.InstanceId.Hex())
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

func (ru *Runner) runPitbull(job *models.PitbullJob) {
	logger := common.NewLogger("Runner", os.Stdout, os.Stderr, "runPitbull", job.ID.Hex())

	logger.Info.Printf("starting Pitbull\n")

	if _, err := ru.manager.RunPitbull(job.InstanceId.Hex()); err != nil {
		logger.Err.Printf("starting Pitbull failed. Reason: %s\n", err)

		ru.handleFailure(job)
		return
	}

	ticker := time.NewTicker(checkPitbullInterval)

	retryCount := 0

	lastProgress := int64(0)
	stalledProgressCount := 0

	for range ticker.C {
		if retryCount >= checkStatusRetryLimit {
			logger.Err.Printf("retries limit reached, stopping job and host\n")

			if err := ru.manager.StopHostInstance(job.InstanceId.Hex()); err != nil {
				logger.Err.Printf("stopping host instance failed. Reason: %s\n", err)
			}

			ticker.Stop()

			ru.handleFailure(job)
			return
		}

		instance, err := ru.manager.SyncInstance(job.InstanceId.Hex())
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

			output, err := ru.manager.GetInstanceOutput(instance)
			if err != nil {
				logger.Err.Printf("output retrieval failed. Reason: %s\n", err)
			}

			if err := ru.manager.StopHostInstance(job.InstanceId.Hex()); err != nil {
				logger.Err.Printf("stopping host instance '%d' failed. reason: %s\n", instance.HostInstance.ProviderId(), err)
			}

			logger.Info.Printf("host instance stopped \n")

			if output != "" {
				instance.LastOutput = output

				if err := ru.manager.UpdateInstance(instance); err != nil {
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

			if err := ru.manager.StopHostInstance(job.InstanceId.Hex()); err != nil {
				logger.Err.Printf("stopping host instance failed. Reason: %s\n", err)
			}

			ticker.Stop()

			ru.handleFailure(job)
			return
		}

		lastProgress = instance.Progress.Checked

	}
}

// handleCompletion - performs any cleanups and updates after completing the job.
func (ru *Runner) handleCompletion(job *models.PitbullJob) {
	logger := common.NewLogger("Runner", os.Stdout, os.Stderr, "cleanup", job.ID.Hex())

	if err := ru.queue.Ack(job.ID.Hex()); err != nil {
		logger.Err.Printf("Acknowledge failed. reason: %s\n", err)
		return
	}

	logger.Info.Printf("Job acknowledged\n")

	job.Status = models.Acknowledged
	job.AcknowledgedAt = time.Now()

	if err := ru.jobRepository.Update(job); err != nil {
		logger.Err.Printf("Updating status failed. reason: %s\n", err)
	}

	return
}

// handleFailure - handles a failure scenario by rescheduling or rejecting the job,
// based on its history and status.
func (ru *Runner) handleFailure(job *models.PitbullJob) {
	logger := common.NewLogger("Runner", os.Stdout, os.Stderr, "cleanup", job.ID.Hex())

	logger.Info.Printf("Failure handling started\n")

	if job.RescheduleCount > rescheduleLimit {
		logger.Info.Printf("Reschedule limit reached, rejecting\n")

		if err := ru.queue.Reject(job.ID.Hex()); err != nil {
			logger.Err.Printf("Rejecting failed. reason: %s\n", err)
			return
		}

		job.Status = models.Rejected
		job.RejectedAt = time.Now()
	} else {
		logger.Info.Printf("Rescheduling.\n")

		if err := ru.queue.Reschedule(job.ID.Hex()); err != nil {
			logger.Err.Printf("Rescheduling failed. reason: %s\n", err)
			return
		}

		job.Status = models.Rescheduled
		job.LastScheduledAt = time.Now()
		job.RescheduleCount += 1
	}

	if err := ru.jobRepository.Update(job); err != nil {
		logger.Info.Printf("Updating job failed. reason: %s\n", err)
		return
	}
}
