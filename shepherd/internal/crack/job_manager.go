package crack

import (
	"fmt"
	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/generator"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull"
	"os"
)

// JobManager - simple facade for operations on PitbullJobs.
type JobManager struct {
	instanceManager *pitbull.InstanceManager
	jobRepository   *JobRepository
	jobQueue        *JobQueue
	tokenGenerator  *generator.TokenGenerator

	logger *common.Logger
}

// NewJobManager - returns new JobService instance.
func NewJobManager(instanceManager *pitbull.InstanceManager) *JobManager {
	return &JobManager{
		instanceManager: instanceManager,
		jobRepository:   NewJobRepository(),
		jobQueue:        NewJobQueue(),
		tokenGenerator:  generator.NewTokenGenerator(generator.TokenRulesetOne),

		logger: common.NewLogger("JobManager", os.Stdout, os.Stderr),
	}
}

func (jm *JobManager) ScheduleJob(job *models.CrackJob) error {
	if err := jm.jobQueue.Enqueue(job.ID.Hex()); err != nil {
		return err
	}

	jm.logger.Info.Printf("job '%s' scheduled\n", job.ID.Hex())

	return nil
}

// GetJobs - returns all existing jobs.
func (jm *JobManager) GetJobs(payload *models.CrackJobsListPayload) ([]*models.CrackJob, int, error) {
	return jm.jobRepository.GetAll(payload)
}

// GetJob - returns a single CrackJob with provided ID.
func (jm *JobManager) GetJob(jobId string) (*models.CrackJob, error) {
	return jm.jobRepository.GetById(jobId)
}

func (jm *JobManager) CreateJob(walletString string, payload *models.CrackPayload, scheduleRun bool) (*models.CrackJob, error) {
	job := models.NewCrackJob(walletString)

	job.FirstScheduledAt = models.NullableTimeNow()
	job.LastScheduledAt = models.NullableTimeNow()

	if payload.Name != "" {
		job.Name = payload.Name
	}

	if len(payload.Tokens) > 0 {
		job.Tokens = payload.Tokens
	}

	if job.Tokens == nil && payload.Keyword != "" {
		job.Keyword = payload.Keyword
		job.Tokens = jm.tokenGenerator.Generate(payload.Keyword)
	}

	if job.Tokens == nil && job.Keyword == "" {
		job.PasslistUrl = payload.PasslistUrl
	}

	if err := jm.jobRepository.Create(job); err != nil {
		return nil, err
	}

	if scheduleRun {
		if err := jm.ScheduleJob(job); err != nil {
			return nil, err
		}
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
func (jm *JobManager) RejectJob(jobId string, reason error) error {
	job, err := jm.GetJob(jobId)
	if err != nil {
		return err
	}

	if job.IsFinished() {
		return fmt.Errorf("job is already finished and cannot be rejected")
	}

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

	return jm.RejectJob(job.ID.Hex(), fmt.Errorf("job has been canceled by the user"))
}

// RescheduleJob - reschedules a single job and marks related instances as "Failed".
func (jm *JobManager) RescheduleJob(jobId string, reason error) error {
	job, err := jm.GetJob(jobId)
	if err != nil {
		return err
	}

	if job.IsFinished() {
		return fmt.Errorf("job is already finished and cannot be rescheduled")
	}

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
	// so any dangling Pitbull instances can be picked up and destroyed
	// by InstanceCollector.
	if err := jm.jobRepository.RescheduleProcessingJobs(jobIds); err != nil {
		return nil, err
	}

	return jobIds, nil
}

// RecreateJob - recreates already finished CrackJob with given ID.
func (jm *JobManager) RecreateJob(jobId string, scheduleRun bool) (*models.CrackJob, error) {
	job, err := jm.GetJob(jobId)
	if err != nil {
		return nil, err
	}

	if !job.IsFinished() {
		return nil, fmt.Errorf("job is not finished and cannot be recreated")
	}

	// Even though some of those fields may be empty,
	// we simply want to recreate the payload used for the original job.
	payload := &models.CrackPayload{
		Name:        job.Name,
		PasslistUrl: job.PasslistUrl,
	}

	// If original's job Keyword is not empty, it means it was used to create the job.
	// Otherwise, Tokens were used directly.
	if job.Keyword != "" {
		payload.Keyword = job.Keyword
	} else {
		payload.Tokens = job.Tokens
	}

	return jm.CreateJob(job.WalletString, payload, scheduleRun)
}
