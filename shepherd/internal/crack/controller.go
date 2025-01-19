package crack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/common/api"
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull"
)

// Controller - responsible for handling all endpoints related to Pitbull
// instances and jobs.
type Controller struct {
	responseHelper api.ResponseHelper
	logger         *common.Logger

	appConfig *config.AppConfig

	jobManager *JobManager
}

// NewController - returns new Controller instance.
func NewController() *Controller {
	instanceManager := pitbull.NewInstanceManager()
	logger := common.NewLogger("Controller", os.Stdout, os.Stderr)

	return &Controller{
		responseHelper: *api.NewResponseHelper(logger),
		appConfig:      config.GetAppConfig(),

		jobManager: NewJobManager(instanceManager),

		logger: logger,
	}
}

// Crack - runs a single CrackJob, based on given CrackPayload.
func (ct *Controller) Crack(
	w http.ResponseWriter,
	r *http.Request,
) {
	var payload *models.CrackPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, err)
		return
	}

	if len(payload.Keywords) == 0 && payload.PasslistUrl == "" && payload.Tokenlist == "" {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, fmt.Errorf("keyword, passlistUrl or tokenlist must be provided"))
		return
	}

	jobs, err := ct.jobManager.HandleJobCreation(ct.appConfig.WalletString, payload, true)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(jobs)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.responseHelper.HandleJSONResponse(w, response)
}

// GetJobs - returns PitbullJobs based on passed filters.
func (ct *Controller) GetJobs(
	w http.ResponseWriter,
	r *http.Request,
) {
	payload := models.NewCrackJobsListPayload()

	if err := payload.Populate(r); err != nil {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, err)
		return
	}

	jobs, totalCount, err := ct.jobManager.GetJobs(payload)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	jobsRaw, err := json.Marshal(jobs)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(api.NewListResponse(jobsRaw, payload.Page, totalCount))
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.responseHelper.HandleJSONResponse(w, response)
}

// CancelJob - rejects a single CrackJob.
func (ct *Controller) CancelJob(
	w http.ResponseWriter,
	r *http.Request,
) {
	var payload *models.CancelCrackJobPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, err)
		return
	}

	if payload.JobId == "" {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, fmt.Errorf("jobId was not provided"))
		return
	}

	job, err := ct.jobManager.GetJob(payload.JobId)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	if err := ct.jobManager.CancelJob(job); err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	updatedJob, err := ct.jobManager.GetJob(payload.JobId)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(updatedJob)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.responseHelper.HandleJSONResponse(w, response)
}

func (ct *Controller) RecreateJob(
	w http.ResponseWriter,
	r *http.Request,
) {
	var payload *models.RecreateCrackJobPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, err)
		return
	}

	if payload.JobId == "" {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, fmt.Errorf("jobId was not provided"))
		return
	}

	job, err := ct.jobManager.RecreateJob(payload.JobId, true)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(job)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.responseHelper.HandleJSONResponse(w, response)
}
