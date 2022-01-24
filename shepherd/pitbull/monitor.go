package pitbull

import (
	"log"
	"os"

	"github.com/el-mike/dogecrack/shepherd/models"
	"github.com/robfig/cron/v3"
)

// PitbullMonitor - entity responsible for monitoring and maintanence jobs.
type PitbullMonitor struct {
	pitbullManager *PitbullManager
	errorLogger    *log.Logger
	infoLogger     *log.Logger
}

func NewPitbullMonitor(pitbullManager *PitbullManager) *PitbullMonitor {
	return &PitbullMonitor{
		pitbullManager: pitbullManager,

		errorLogger: log.New(os.Stderr, "[Monitor][Error]: ", log.Ldate|log.Ltime),
		infoLogger:  log.New(os.Stdout, "[Monitor][Info]: ", log.Ldate|log.Ltime),
	}
}

func (pm *PitbullMonitor) RunMonitoring(instanceId string) {
	pm.infoLogger.Println("---Monitor: <start>: " + instanceId + " ---")

	c := cron.New()

	c.AddFunc("* * * * *", func() {
		if _, err := pm.monitorJob(instanceId); err != nil {
			pm.errorLogger.Println(err)
			c.Stop()
		}

		pm.infoLogger.Println("---Monitor: <end>---")
	})

	c.Start()
}

func (pm *PitbullMonitor) monitorJob(instanceId string) (*models.PitbullInstance, error) {
	instance, err := pm.pitbullManager.GetInstanceById(instanceId)
	if err != nil {
		return nil, err
	}

	pm.infoLogger.Println(instance.Status)
	pm.infoLogger.Printf("%d / %d\n", instance.Progress.Checked, instance.Progress.Total)

	return instance, nil
}
