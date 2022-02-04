package pitbull

import (
	"errors"
	"os"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/models"
)

const (
	startHostAttemptsLimit = 10
	checkStatusRetryLimit  = 10
	checkHostInterval      = 15 * time.Second
	checkPitbullInterval   = 30 * time.Second
)

// Runner - entity responsible for running and monitoring Pitbull jobs.
type Runner struct {
	manager *Manager
}

// NewRunner - returns new PitbullRunner instance.
func NewRunner(manager *Manager) *Runner {
	return &Runner{
		manager: manager,
	}
}

// Run - starts single Pitbull run.
func (ru *Runner) Run(instanceId string) {
	go ru.startHost(instanceId)
}

func (ru *Runner) startHost(instanceId string) {
	logger := common.NewLogger("Runner", os.Stdout, os.Stderr, "startHost", instanceId)

	logger.Info.Println("starting host.")

	_, err := ru.manager.RunHostForInstance(instanceId)
	if err != nil {
		logger.Err.Printf("starting host failed. Reason: %s\n", err)

		return
	}

	ticker := time.NewTicker(checkHostInterval)

	attemptsCount := 0

	for range ticker.C {
		attemptsCount += 1

		if attemptsCount >= startHostAttemptsLimit {
			logger.Info.Printf("attempts limit reached, stopping job and host\n")

			if err := ru.manager.StopHostInstance(instanceId); err != nil {
				logger.Err.Printf("stopping host instance failed. Reason: %s\n", err)
			}

			ticker.Stop()
			return
		}

		instance, err := ru.manager.SyncInstance(instanceId)
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

			go ru.runPitbull(instanceId)

			ticker.Stop()
			break
		}
	}
}

func (ru *Runner) runPitbull(instanceId string) {
	logger := common.NewLogger("Runner", os.Stdout, os.Stderr, "runPitbullJob", instanceId)

	logger.Info.Printf("starting Pitbull\n")

	if _, err := ru.manager.RunPitbull(instanceId); err != nil {
		logger.Err.Printf("starting Pitbull failed. Reason: %s\n", err)

		return
	}

	ticker := time.NewTicker(checkPitbullInterval)

	retryCount := 0

	for range ticker.C {
		if retryCount >= checkStatusRetryLimit {
			logger.Info.Printf("retries limit reached, stopping job and host\n")

			if err := ru.manager.StopHostInstance(instanceId); err != nil {
				logger.Err.Printf("stopping host instance failed. Reason: %s\n", err)
			}

			ticker.Stop()
			return
		}

		instance, err := ru.manager.SyncInstance(instanceId)
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

			if err := ru.manager.StopHostInstance(instanceId); err != nil {
				logger.Err.Printf("stopping host instance '%d' failed. reason: %s\n", instance.HostInstance.ProviderId(), err)
			}

			logger.Info.Printf("host instance stopped \n")

			if output != "" {
				instance.LastOutput = output

				if err := ru.manager.UpdateInstance(instance); err != nil {
					logger.Err.Printf("saving last output failed. Reason: %s\n", err)
				}
			}

			ticker.Stop()

			logger.Info.Printf("job completed\n")
			return
		}
	}
}
