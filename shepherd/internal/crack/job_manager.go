package crack

import (
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/generator"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull"
)

// JobManager - simple facade for operations on PitbullJobs.
type JobManager struct {
	jobRepository   *JobRepository
	instanceManager *pitbull.InstanceManager
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
func (jm *JobManager) GetJobs(payload *models.PitbullJobsListPayload) ([]*models.CrackJob, int, error) {
	return jm.jobRepository.GetAll(payload)
}

// CreateKeywordJob - creates new CrackJob based on keyword, with generated tokens.
func (jm *JobManager) CreateKeywordJob(walletString, keyword string) (*models.CrackJob, error) {
	job := jm.createBaseJob(walletString)

	tokens := jm.tokenGenerator.Generate(keyword)
	job.Keyword = keyword
	job.Tokens = tokens

	if err := jm.jobRepository.Create(job); err != nil {
		return nil, err
	}

	return job, nil
}

// CreatePasslistJob - create new CrackJob based on passlist URL.
func (jm *JobManager) CreatePasslistJob(walletString, passlistUrl string) (*models.CrackJob, error) {
	job := jm.createBaseJob(walletString)

	job.PasslistUrl = passlistUrl

	if err := jm.jobRepository.Create(job); err != nil {
		return nil, err
	}

	return job, nil
}

func (jm *JobManager) createBaseJob(walletString string) *models.CrackJob {
	job := models.NewCrackJob(walletString)

	job.FirstScheduledAt = models.NullableTimeNow()
	job.LastScheduledAt = models.NullableTimeNow()

	return job
}

// AssignInstance - creates a PitbullInstance and assigns it to passed CrackJob.
// If job has been rescheduled, it will attempt to destroy previous HostInstance.
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
