package pitbull

import (
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/models"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/repositories"
)

// JobManager - simple facade for operations on PitbullJobs.
type JobManager struct {
	jobRepository   *repositories.JobRepository
	instanceManager *InstanceManager
	jobQueue        *JobQueue
}

// NewJobManager - returns new JobService instance.
func NewJobManager(instanceManager *InstanceManager) *JobManager {
	return &JobManager{
		jobRepository:   repositories.NewJobRepository(),
		instanceManager: instanceManager,
		jobQueue:        NewJobQueue(),
	}
}

// GetJobs - returns all existing jobs.
func (js *JobManager) GetJobs(statuses []models.JobStatus) ([]*models.PitbullJob, error) {
	return js.jobRepository.GetAll(statuses)
}

// CreateJob - creates an empty job and saves it to DB.
func (js *JobManager) CreateJob(passlistUrl, walletString string) (*models.PitbullJob, error) {
	job := models.NewPitbullJob(passlistUrl, walletString)

	job.FirstScheduledAt = time.Now()
	job.LastScheduledAt = time.Now()

	if err := js.jobRepository.Create(job); err != nil {
		return nil, err
	}

	return job, nil
}

// AssignInstance - creates a PitbullInstance and assignes it to passed PitbullJob.
// If job has been rescheduled, it will attempt to destroy previous HostInstance.
func (js *JobManager) AssignInstance(job *models.PitbullJob) (*models.PitbullJob, error) {
	instance, err := js.instanceManager.CreateInstance(job.PasslistUrl, job.WalletString)
	if err != nil {
		return nil, err
	}

	job.InstanceId = instance.ID

	if err := js.jobRepository.Update(job); err != nil {
		return nil, err
	}

	return job, nil
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

	job.Status = models.JobProcessing
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

	job.Status = models.JobAcknowledged
	job.AcknowledgedAt = time.Now()

	return js.jobRepository.Update(job)
}

// RejectJob - rejects a single job, and marks related instances as "Interrupted".
func (js *JobManager) RejectJob(job *models.PitbullJob) error {
	if err := js.jobQueue.Reject(job.ID.Hex()); err != nil {
		return err
	}

	job.Status = models.JobRejected
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

	job.Status = models.JobRescheduled
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
	if err := js.jobRepository.RescheduleProcessingJobs(jobIds); err != nil {
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

	return js.instanceManager.UpdateInstance(instance)
}
