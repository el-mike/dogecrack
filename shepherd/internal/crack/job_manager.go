package crack

import (
	"fmt"
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/generator"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull"
)

// JobManager - simple facade for operations on PitbullJobs.
type JobManager struct {
	instanceManager *pitbull.InstanceManager
	jobRepository   *JobRepository
	jobQueue        *JobQueue
	tokenGenerator  *generator.TokenGenerator
}

// NewJobManager - returns new JobService instance.
func NewJobManager(instanceManager *pitbull.InstanceManager) *JobManager {
	return &JobManager{
		instanceManager: instanceManager,
		jobRepository:   NewJobRepository(),
		jobQueue:        NewJobQueue(),
		tokenGenerator:  generator.NewTokenGenerator(generator.TokenRulesetOne),
	}
}

// GetJobs - returns all existing jobs.
func (jm *JobManager) GetJobs(payload *models.CrackJobsListPayload) ([]*models.CrackJob, int, error) {
	return jm.jobRepository.GetAll(payload)
}

// GetJob - returns a single CrackJob with provided ID.
func (jm *JobManager) GetJob(jobId string) (*models.CrackJob, error) {
	return jm.jobRepository.GetById(jobId)
}

func (jm *JobManager) CreateJob(walletString string, payload *models.CrackPayload) (*models.CrackJob, error) {
	job := models.NewCrackJob(walletString)

	job.FirstScheduledAt = models.NullableTimeNow()
	job.LastScheduledAt = models.NullableTimeNow()

	if payload.Name != "" {
		job.Name = payload.Name
	}

	if payload.Keyword != "" {
		job.Keyword = payload.Keyword
		job.Tokens = jm.tokenGenerator.Generate(payload.Keyword)
	} else {
		job.PasslistUrl = payload.PasslistUrl
	}

	if err := jm.jobRepository.Create(job); err != nil {
		return nil, err
	}

	return job, nil
}

// AssignInstance - creates a PitbullInstance and assigns it to passed CrackJob.
func (jm *JobManager) AssignInstance(job *models.CrackJob) (*models.CrackJob, error) {
	runPayload := &models.PitbullRunPayload{
		WalletString: job.WalletString,
		Tokenlist:    job.GetTokenlist(),
		PasslistUrl:  job.PasslistUrl,
	}

	instance, err := jm.instanceManager.CreateInstance(job.ID, runPayload)
	if err != nil {
		return nil, err
	}

	job.InstanceId = instance.ID

	if err := jm.jobRepository.Update(job); err != nil {
		return nil, err
	}

	return job, nil
}

// DequeueJob - dequeues a job and marks it as "Processing".
func (jm *JobManager) DequeueJob() (*models.CrackJob, error) {
	jobId, err := jm.jobQueue.Dequeue()
	if err != nil {
		return nil, err
	}

	// Empty jobId should not be treated as error - it just means that the
	// "workingQueue" is empty.
	if jobId == "" {
		return nil, nil
	}

	job, err := jm.jobRepository.GetById(jobId)
	if err != nil {
		return nil, err
	}

	job.Status = models.JobStatus.Processing
	job.StartedAt = models.NullableTimeNow()

	if err := jm.jobRepository.Update(job); err != nil {
		return nil, err
	}

	return job, nil
}

// AcknowledgeJob - ackes a single job.
func (jm *JobManager) AcknowledgeJob(job *models.CrackJob) error {
	if err := jm.jobQueue.Ack(job.ID.Hex()); err != nil {
		return err
	}

	job.Status = models.JobStatus.Acknowledged
	job.AcknowledgedAt = models.NullableTimeNow()

	return jm.jobRepository.Update(job)
}

// RejectJob - rejects a single job, and marks related instances as "Failed".
func (jm *JobManager) RejectJob(job *models.CrackJob, reason error) error {
	if err := jm.jobQueue.Reject(job.ID.Hex()); err != nil {
		return err
	}

	job.AppendError(reason)

	job.Status = models.JobStatus.Rejected
	job.RejectedAt = models.NullableTimeNow()

	if err := jm.jobRepository.Update(job); err != nil {
		return err
	}

	return jm.instanceManager.MarkInstanceAsFailed(job.InstanceId.Hex(), reason)
}

// CancelJob - destroys currently assigned instance and rejects a job.
func (jm *JobManager) CancelJob(job *models.CrackJob) error {
	if job.IsFinished() {
		return fmt.Errorf("job is already finished and cannot be canceled")
	}

	err := jm.instanceManager.StopHostInstance(job.InstanceId.Hex())
	if err != nil {
		return err
	}

	return jm.RejectJob(job, fmt.Errorf("job has been canceled by the user"))
}

// RescheduleJob - reschedules a single job and marks related instances as "Failed".
func (jm *JobManager) RescheduleJob(job *models.CrackJob, reason error) error {
	if err := jm.jobQueue.Reschedule(job.ID.Hex()); err != nil {
		return err
	}

	job.AppendError(reason)

	job.Status = models.JobStatus.Rescheduled
	job.LastScheduledAt = models.NullableTimeNow()
	job.RescheduleCount += 1

	if err := jm.jobRepository.Update(job); err != nil {
		return err
	}

	return jm.instanceManager.MarkInstanceAsFailed(job.InstanceId.Hex(), reason)
}

// RescheduleProcessingJobs - reschedules all jobs in "processingQueue".
func (jm *JobManager) RescheduleProcessingJobs() ([]string, error) {
	jobIds, err := jm.jobQueue.RescheduleAllProcessing()
	if err != nil {
		return nil, err
	}

	// We need to reject all previously processing jobs,
	// so any dangling Pitbull instanes can be picked up and destroyed
	// by InstanceCollector.
	if err := jm.jobRepository.RescheduleProcessingJobs(jobIds); err != nil {
		return nil, err
	}

	return jobIds, nil
}
