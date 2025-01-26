package crack

import (
	"fmt"
	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/common/repositories"
	"github.com/el-mike/dogecrack/shepherd/internal/generator"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull"
	"os"
	"strings"
	"sync"
)

// JobManager - simple facade for operations on PitbullJobs.
type JobManager struct {
	appSettingsRepository *repositories.AppSettingsRepository
	jobRepository         *JobRepository
	instanceManager       *pitbull.InstanceManager
	jobQueue              *JobQueue
	tokenGeneratorFactory *generator.TokenGeneratorFactory

	logger *common.Logger
}

// NewJobManager - returns new JobService instance.
func NewJobManager(instanceManager *pitbull.InstanceManager) *JobManager {
	return &JobManager{
		appSettingsRepository: repositories.NewAppSettingsRepository(),
		jobRepository:         NewJobRepository(),
		instanceManager:       instanceManager,
		jobQueue:              NewJobQueue(),
		tokenGeneratorFactory: generator.NewTokenGeneratorFactory(),

		logger: common.NewLogger("JobManager", os.Stdout, os.Stderr),
	}
}

func (jm *JobManager) ScheduleJobs(jobs []*models.CrackJob) error {
	var jobIds []string

	for _, job := range jobs {
		jobIds = append(jobIds, job.ID.Hex())
	}

	if err := jm.jobQueue.Enqueue(jobIds...); err != nil {
		return err
	}

	jm.logger.Info.Printf("jobs: '%s' scheduled\n", strings.Join(jobIds, ", "))

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

func (jm *JobManager) getBaseJob(walletString, name string) *models.CrackJob {
	job := models.NewCrackJob(walletString)

	job.FirstScheduledAt = models.NullableTimeNow()
	job.LastScheduledAt = models.NullableTimeNow()

	if name != "" {
		job.Name = name
	}

	return job
}

func (jm *JobManager) CreateJobsForKeywords(
	walletString,
	name string,
	keywords []string,
	tokenGeneratorVersion models.TokenGeneratorVersionEnum,
	scheduleRun bool,
) ([]*models.CrackJob, error) {
	var jobs []*models.CrackJob

	for _, keyword := range keywords {
		job := jm.getBaseJob(walletString, name)

		generatorVersion := models.LatestTokenGeneratorVersion

		if tokenGeneratorVersion != 0 {
			generatorVersion = tokenGeneratorVersion
		}

		tokenGenerator, err := jm.tokenGeneratorFactory.CreateGenerator(generatorVersion)
		if err != nil {
			return nil, err
		}

		job.Keyword = keyword
		job.Tokenlist = tokenGenerator.Generate(keyword)
		job.TokenGeneratorVersion = generatorVersion

		jobs = append(jobs, job)
	}

	if err := jm.jobRepository.CreateMany(jobs); err != nil {
		return nil, err
	}

	if scheduleRun {
		if err := jm.ScheduleJobs(jobs); err != nil {
			return nil, err
		}
	}

	return jobs, nil
}

func (jm *JobManager) CreateJobForTokenlist(
	walletString,
	name,
	tokenlist,
	baseKeyword string,
	scheduleRun bool,
) (*models.CrackJob, error) {
	job := jm.getBaseJob(walletString, name)

	job.CustomTokenlist = true
	job.Tokenlist = tokenlist

	// If base keyword was provided, set it as Keyword attribute.
	if baseKeyword != "" {
		job.Keyword = baseKeyword
	}

	if err := jm.jobRepository.Create(job); err != nil {
		return nil, err
	}

	if scheduleRun {
		if err := jm.ScheduleJobs([]*models.CrackJob{job}); err != nil {
			return nil, err
		}
	}

	return job, nil
}

func (jm *JobManager) CreateJobForPasslist(walletString, name, passlistUrl string, scheduleRun bool) (*models.CrackJob, error) {
	job := jm.getBaseJob(walletString, name)

	job.PasslistUrl = passlistUrl

	if err := jm.jobRepository.Create(job); err != nil {
		return nil, err
	}

	if scheduleRun {
		if err := jm.ScheduleJobs([]*models.CrackJob{job}); err != nil {
			return nil, err
		}
	}

	return job, nil
}

func (jm *JobManager) HandleJobCreation(walletString string, payload *models.CrackPayload, scheduleRun bool) ([]*models.CrackJob, error) {
	// If custom tokenlist was provided, we want to prioritize it.
	if payload.Tokenlist != "" {
		keyword := ""

		// If base keyword was provided, we want to save this information as well.
		if len(payload.Keywords) > 0 {
			keyword = payload.Keywords[0]
		}

		job, err := jm.CreateJobForTokenlist(walletString, payload.Name, payload.Tokenlist, keyword, scheduleRun)
		if err != nil {
			return nil, err
		}

		return []*models.CrackJob{job}, nil
	}

	if len(payload.Keywords) > 0 {
		return jm.CreateJobsForKeywords(walletString, payload.Name, payload.Keywords, payload.TokenGeneratorVersion, scheduleRun)
	}

	if payload.PasslistUrl != "" {
		job, err := jm.CreateJobForPasslist(walletString, payload.Name, payload.PasslistUrl, scheduleRun)
		if err != nil {
			return nil, err
		}

		return []*models.CrackJob{job}, nil
	}

	return nil, fmt.Errorf("no keywords, passlistUrl or tokenlist provided")
}

// AssignInstance - creates a PitbullInstance and assigns it to passed CrackJob.
func (jm *JobManager) AssignInstance(job *models.CrackJob, previousInstance *models.PitbullInstance) (*models.CrackJob, error) {
	appSettings, err := jm.appSettingsRepository.GetAppSettings()
	if err != nil {
		return nil, err
	}

	runPayload := &models.PitbullRunPayload{
		WalletString: job.WalletString,
		Tokenlist:    job.Tokenlist,
		PasslistUrl:  job.PasslistUrl,
	}

	if appSettings.MinPasswordLength > 0 {
		runPayload.MinLength = appSettings.MinPasswordLength
	}

	if appSettings.MaxPasswordLength > 0 {
		runPayload.MaxLength = appSettings.MaxPasswordLength
	}

	if previousInstance != nil {
		runPayload.SkipCount = previousInstance.Pitbull.GetResumeCount()
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

// GetProcessingJobsIDs - returns  IDs of CrackJobs currently being run.
func (jm *JobManager) GetProcessingJobsIDs() ([]string, error) {
	return jm.jobQueue.GetProcessing()
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
		Name:                  job.Name,
		TokenGeneratorVersion: job.TokenGeneratorVersion,
	}

	// If PasslistUrl is not empty, it means it was used to create the job.
	// Otherwise, if TokenGeneratorVersion is equal to zero, we know that custom tokenlist was provided,
	// so we set it explicitly in payload as well (which will make the JobManager skip token generation).
	if job.PasslistUrl != "" {
		payload.PasslistUrl = job.PasslistUrl
	} else if job.TokenGeneratorVersion == 0 {
		payload.Tokenlist = job.Tokenlist
	}

	// This covers both keyword and custom tokenlist usage.
	if job.Keyword != "" {
		payload.Keywords = []string{job.Keyword}
	}

	result, err := jm.HandleJobCreation(job.WalletString, payload, scheduleRun)
	if err != nil {
		return nil, err
	}

	return result[0], nil
}

// GetCheckedIdeas - returns CheckedKeywords and CheckedPasswords.
func (jm *JobManager) GetCheckedIdeas() (*models.CheckedIdeas, error) {
	checkedIdeas := &models.CheckedIdeas{}

	var mainErr error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		keywords, err := jm.jobRepository.GetCheckedKeywords()
		if err != nil {
			mainErr = err
			return
		}

		checkedIdeas.CheckedKeywords = keywords
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		passlists, err := jm.jobRepository.GetCheckedPasslists()
		if err != nil {
			mainErr = err
			return
		}

		checkedIdeas.CheckedPasslists = passlists
	}()

	wg.Wait()

	if mainErr != nil {
		return nil, mainErr
	}

	return checkedIdeas, nil
}
