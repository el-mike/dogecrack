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
	jobQueue           *JobQueue
}

// NewJobManager - returns new JobService instance.
func NewJobManager() *JobManager {
	return &JobManager{
		jobRepository:      repositories.NewJobRepository(),
		instanceRepository: repositories.NewInstanceRepository(),
		jobQueue:           NewJobQueue(),
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

// DequeueJob - dequeues a job and marks it as "Processing".
func (js *JobManager) DequeueJob() (*models.PitbullJob, error) {
	jobId, err := js.jobQueue.Dequeue()
	if err != nil {
		return nil, err
	}

	// Empty jobId should not be treated as error - it just means that the
	// "workingQueue" is empty.
	if jobId == "" {
		return nil, nil
	}

	job, err := js.jobRepository.GetById(jobId)
	if err != nil {
		return nil, err
	}

	job.Status = models.Processing
	job.StartedAt = time.Now()

	if err := js.jobRepository.Update(job); err != nil {
		return nil, err
	}

	return job, nil
}

// AcknowledgeJob - ackes a single job.
func (js *JobManager) AcknowledgeJob(job *models.PitbullJob) error {
	if err := js.jobQueue.Ack(job.ID.Hex()); err != nil {
		return err
	}

	job.Status = models.Acknowledged
	job.AcknowledgedAt = time.Now()

	return js.jobRepository.Update(job)
}

// RejectJob - rejects a single job, and marks related instances as "Interrupted".
func (js *JobManager) RejectJob(job *models.PitbullJob) error {
	if err := js.jobQueue.Reject(job.ID.Hex()); err != nil {
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
	if err := js.jobQueue.Reschedule(job.ID.Hex()); err != nil {
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

// RescheduleProcessingJobs - reschedules all jobs in "processingQueue".
func (js *JobManager) RescheduleProcessingJobs() ([]string, error) {
	jobIds, err := js.jobQueue.RescheduleAllProcessing()
	if err != nil {
		return nil, err
	}

	// We need to reject all previously processing jobs,
	// so any dangling Pitbull instanes can be picked up and destroyed
	// by InstanceCollector.
	if err := js.jobRepository.RejectProcessingJobs(jobIds); err != nil {
		return nil, err
	}

	return jobIds, nil
}

// MarkInstanceAsInterrupted - marks given instance as "Interrupted" and updates it
// in the DB.
func (js *JobManager) MarkInstanceAsInterrupted(instance *models.PitbullInstance) error {
	if instance == nil {
		return nil
	}

	instance.Status = models.Interrupted

	return js.instanceRepository.UpdateInstance(instance)
}
