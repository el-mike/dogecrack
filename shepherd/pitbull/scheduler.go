package pitbull

import (
	"log"
	"os"
	"sync"

	"github.com/el-mike/dogecrack/shepherd/models"
	"github.com/robfig/cron/v3"
)

const INSTANCES_LIMIT = 5
const START_HOST_ATTEMPTS_LIMIT = 10

// PitbullScheduler - entity responsible for monitoring and maintanence jobs.
type PitbullScheduler struct {
	pitbullManager *PitbullManager
	errorLogger    *log.Logger
	infoLogger     *log.Logger

	instancesLimit int
	instancesQueue *models.StringQueue

	jobsPool *jobsPool
}

func NewPitbullScheduler(pitbullManager *PitbullManager) *PitbullScheduler {
	return &PitbullScheduler{
		pitbullManager: pitbullManager,

		errorLogger: log.New(os.Stderr, "[Scheduler][Error] ", log.Ldate|log.Ltime),
		infoLogger:  log.New(os.Stdout, "[Scheduler][Info] ", log.Ldate|log.Ltime),

		instancesLimit: INSTANCES_LIMIT,
		instancesQueue: models.NewStringQueue(),

		jobsPool: newJobsPool(),
	}
}

// ScheduleRun - schedules a single Pitbull run. If instances limit is not reach yet,
// it will run it immediately.
func (pm *PitbullScheduler) ScheduleRun(passlistUrl, walletString string) (*models.PitbullInstance, error) {
	instance, err := pm.pitbullManager.CreateInstance(passlistUrl, walletString)
	if err != nil {
		return nil, err
	}

	// If the queue is empty, we want to wake the dequeueJob.
	if pm.instancesQueue.IsEmpty() {
		if err := pm.dequeueJob(); err != nil {
			return nil, err
		}
	}

	instanceId := instance.ID.Hex()
	pm.instancesQueue.Enqueue(instanceId)

	pm.infoLogger.Printf("[scheduleRun][%s]: instance scheduled\n", instanceId)

	return instance, nil
}

func (pm *PitbullScheduler) dequeueJob() error {
	pm.infoLogger.Printf("[dequeueJob]: starting\n")

	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("*/30 * * * * *", func() {
		if !pm.instancesQueue.IsEmpty() && pm.jobsPool.Size() < INSTANCES_LIMIT {
			instanceId, err := pm.instancesQueue.Dequeue()
			if err != nil {
				pm.errorLogger.Printf("[dequeueJob][%s]: instance dequeue failed, reason: %s\n", instanceId, err)
			}

			pm.infoLogger.Printf("[dequeueJob][%s]: instance dequeued\n", instanceId)

			pm.startHostJob(instanceId)
		}

		// If queue has been emptied, we want to stop dequeueJob, to not waste resources
		// and prevent memory leaks. It will be waken up by ScheduleRun if needed.
		if pm.instancesQueue.IsEmpty() {
			pm.infoLogger.Printf("[dequeueJob]: queue empty, stopping\n")

			c.Stop()
			return
		}
	})

	if err != nil {
		pm.errorLogger.Printf("[dequeueJob]: cronjob error, reason: %s\n", err)
		return err
	}

	c.Start()

	return nil
}

func (pm *PitbullScheduler) startHostJob(instanceId string) error {
	pm.infoLogger.Printf("[startHostJob][%s]: starting host\n", instanceId)

	_, err := pm.pitbullManager.RunHostForInstance(instanceId)
	if err != nil {
		pm.errorLogger.Printf("[startHostJob][%s]: starting host failed, reason: %s\n", instanceId, err)

		return err
	}

	c := cron.New(cron.WithSeconds())

	attemptsCount := 0

	// Runs at every 30th second.
	_, err = c.AddFunc("*/30 * * * * *", func() {
		instance, err := pm.pitbullManager.SyncInstance(instanceId)
		if err != nil {
			pm.errorLogger.Printf("[startHostJob][%s]: instance sync failed, stopping the job. reason: %s\n", instanceId, err)

			c.Stop()
			return
		}

		pm.infoLogger.Printf("[startHostJob][%s]: host status: %s\n", instanceId, instance.HostInstance.HostStatus().Formatted())

		if instance.Status == models.Waiting {
			pm.infoLogger.Printf("[startHostJob][%s]: host started\n", instanceId)

			pm.runPitbullJob(instanceId)

			c.Stop()
			return
		}

		attemptsCount += 1

		if attemptsCount >= START_HOST_ATTEMPTS_LIMIT {
			pm.infoLogger.Printf("[startHostJob][%s]: attempts limit reached, stopping job and host\n", instanceId)

			if err := pm.pitbullManager.StopHostInstance(instanceId); err != nil {
				pm.errorLogger.Printf("[startHostJob][%s]: stopping host instance failed\n", instanceId)
			}

			c.Stop()
			return
		}
	})

	if err != nil {
		pm.errorLogger.Printf("[startHostJob][%s]: cronjob error, reason: %s\n", instanceId, err)
		return err
	}

	c.Start()

	return nil
}

func (pm *PitbullScheduler) runPitbullJob(instanceId string) error {
	pm.infoLogger.Printf("[runPitbullJob][%s]: starting Pitbull\n", instanceId)

	if _, err := pm.pitbullManager.RunPitbull(instanceId); err != nil {
		pm.infoLogger.Printf("[runPitbullJob][%s]: starting Pitbull failed, reason: %s\n", instanceId, err)

		return err
	}

	// This cron does not use WithSecond(), therefore first place
	// is minutes, not seconds.
	c := cron.New()

	_, err := c.AddFunc("* * * * *", func() {
		instance, err := pm.pitbullManager.SyncInstance(instanceId)
		if err != nil {
			pm.errorLogger.Printf("[runPitbullJob][%s]: Pitbull sync failed, stopping the job. reason: %s\n", instanceId, err)

			c.Stop()
			return
		}

		pm.infoLogger.Printf("[runPitbullJob][%s][Process]: %s | %s\n", instanceId, instance.Status.Formatted(), instance.Progress.Formatted())

		if instance.Status == models.Finished ||
			instance.Status == models.Success {
			pm.infoLogger.Printf("[runPitbullJob][%s]: pitbull finished, stopping host instance\n", instanceId)

			output, err := pm.pitbullManager.GetInstanceOutput(instance)
			if err != nil {
				pm.errorLogger.Printf("[runPitbullJob][%s] output retrieval failed, reason: %s\n", instanceId, err.Error())
			}

			if err := pm.pitbullManager.StopHostInstance(instanceId); err != nil {
				pm.errorLogger.Printf("[runPitbullJob][%s]: stopping host instance '%d' failed\n, reason: %s", instanceId, instance.HostInstance.ProviderId(), err)
			}

			pm.infoLogger.Printf("[runPitbullJob][%s]: host instance stopped \n", instanceId)

			if output != "" {
				instance.LastOutput = output

				if err := pm.pitbullManager.UpdateInstance(instance); err != nil {
					pm.errorLogger.Printf("[runPitbullJob][%s] saving last output failed, reason: %s\n", instanceId, err)
				}
			}

			c.Stop()

			pm.infoLogger.Printf("[runPitbullJob][%s]: job completed \n", instanceId)
			return
		}
	})

	if err != nil {
		pm.errorLogger.Printf("[runPitbullJob][%s]: cronjob error, reason: %s\n", instanceId, err)
		return err
	}

	c.Start()

	return nil
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
