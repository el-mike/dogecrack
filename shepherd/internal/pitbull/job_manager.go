package pitbull

import (
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/models"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/repositories"
)

// JobManager - simple facade for operations on PitbullJobs.
type JobManager struct {
	jobRepository      *repositories.JobRepository
	instanceRepository *repositories.InstanceRepository
	queue              *JobQueue
}

// NewJobManager - returns new JobService instance.
func NewJobManager() *JobManager {
	return &JobManager{
		jobRepository:      repositories.NewJobRepository(),
		instanceRepository: repositories.NewInstanceRepository(),
		queue:              NewJobQueue(),
	}
}

// GetJobs - returns all existing jobs.
func (js *JobManager) GetJobs(statuses []models.JobStatus) ([]*models.PitbullJob, error) {
	return js.jobRepository.GetAll(statuses)
}

// GetCompletedJobWithActiveInstance - returns all completed jobs with active instances.
func (js *JobManager) GetCompletedJobWithActiveInstance() ([]*models.PitbullJob, error) {
	return js.jobRepository.GetCompletedWithActiveInstance()
}

// AcknowledgeJob - ackes a single job.
func (js *JobManager) AcknowledgeJob(job *models.PitbullJob) error {
	if err := js.queue.Ack(job.ID.Hex()); err != nil {
		return err
	}

	job.Status = models.Acknowledged
	job.AcknowledgedAt = time.Now()

	return js.jobRepository.Update(job)
}

// RejectJob - rejects a single job, and marks related instances as "Interrupted".
func (js *JobManager) RejectJob(job *models.PitbullJob) error {
	if err := js.queue.Reject(job.ID.Hex()); err != nil {
		return err
	}

	job.Status = models.Rejected
	job.RejectedAt = time.Now()

	if err := js.jobRepository.Update(job); err != nil {
		return err
	}

	return js.MarkInstanceAsInterrupted(job.Instance)
}

// RescheduleJob - reschedules a single job and marks related instances as "Interrupted".
func (js *JobManager) RescheduleJob(job *models.PitbullJob) error {
	if err := js.queue.Reschedule(job.ID.Hex()); err != nil {
		return err
	}

	job.Status = models.Rescheduled
	job.LastScheduledAt = time.Now()
	job.RescheduleCount += 1

	if err := js.jobRepository.Update(job); err != nil {
		return err
	}

	return js.MarkInstanceAsInterrupted(job.Instance)
}

func (js *JobManager) MarkInstanceAsInterrupted(instance *models.PitbullInstance) error {
	if instance == nil {
		return nil
	}

	instance.Status = models.Interrupted

	return js.instanceRepository.UpdateInstance(instance)
}
