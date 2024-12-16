package crack

import (
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull"
)

// JobManager - simple facade for operations on PitbullJobs.
type JobManager struct {
	jobRepository   *JobRepository
	instanceManager *pitbull.InstanceManager
	jobQueue        *JobQueue
}

// NewJobManager - returns new JobService instance.
func NewJobManager(instanceManager *pitbull.InstanceManager) *JobManager {
	return &JobManager{
		instanceManager: instanceManager,
		jobRepository:   NewJobRepository(),
		jobQueue:        NewJobQueue(),
	}
}

// GetJobs - returns all existing jobs.
func (js *JobManager) GetJobs(payload *models.PitbullJobsListPayload) ([]*models.CrackJob, int, error) {
	return js.jobRepository.GetAll(payload)
}

// CreateJob - creates an empty job and saves it to DB.
func (js *JobManager) CreateJob(keyword, passlistUrl, walletString string) (*models.CrackJob, error) {
	job := models.NewPitbullJob(keyword, passlistUrl, walletString)

	job.FirstScheduledAt = models.NullableTimeNow()
	job.LastScheduledAt = models.NullableTimeNow()

	if err := js.jobRepository.Create(job); err != nil {
		return nil, err
	}

	return job, nil
}

// AssignInstance - creates a PitbullInstance and assigns it to passed PitbullJob.
// If job has been rescheduled, it will attempt to destroy previous HostInstance.
func (js *JobManager) AssignInstance(job *models.CrackJob) (*models.CrackJob, error) {
	instance, err := js.instanceManager.CreateInstance(job.ID, job.PasslistUrl, job.WalletString)
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
func (js *JobManager) DequeueJob() (*models.CrackJob, error) {
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

	job.Status = models.JobStatus.Processing
	job.StartedAt = models.NullableTimeNow()

	if err := js.jobRepository.Update(job); err != nil {
		return nil, err
	}

	return job, nil
}

// AcknowledgeJob - ackes a single job.
func (js *JobManager) AcknowledgeJob(job *models.CrackJob) error {
	if err := js.jobQueue.Ack(job.ID.Hex()); err != nil {
		return err
	}

	job.Status = models.JobStatus.Acknowledged
	job.AcknowledgedAt = models.NullableTimeNow()

	return js.jobRepository.Update(job)
}

// RejectJob - rejects a single job, and marks related instances as "Failed".
func (js *JobManager) RejectJob(job *models.CrackJob, reason error) error {
	if err := js.jobQueue.Reject(job.ID.Hex()); err != nil {
		return err
	}

	job.AppendError(reason)

	job.Status = models.JobStatus.Rejected
	job.RejectedAt = models.NullableTimeNow()

	if err := js.jobRepository.Update(job); err != nil {
		return err
	}

	return js.instanceManager.MarkInstanceAsFailed(job.InstanceId.Hex(), reason)
}

// RescheduleJob - reschedules a single job and marks related instances as "Failed".
func (js *JobManager) RescheduleJob(job *models.CrackJob, reason error) error {
	if err := js.jobQueue.Reschedule(job.ID.Hex()); err != nil {
		return err
	}

	job.AppendError(reason)

	job.Status = models.JobStatus.Rescheduled
	job.LastScheduledAt = models.NullableTimeNow()
	job.RescheduleCount += 1

	if err := js.jobRepository.Update(job); err != nil {
		return err
	}

	return js.instanceManager.MarkInstanceAsFailed(job.InstanceId.Hex(), reason)
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
