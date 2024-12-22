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

	jobScheduler *Scheduler
	jobManager   *JobManager
}

// NewController - returns new Controller instance.
func NewController() *Controller {
	instanceManager := pitbull.NewInstanceManager()
	logger := common.NewLogger("Controller", os.Stdout, os.Stderr)

	return &Controller{
		responseHelper: *api.NewResponseHelper(logger),
		appConfig:      config.GetAppConfig(),

		jobScheduler: NewScheduler(instanceManager),
		jobManager:   NewJobManager(instanceManager),

		logger: logger,
	}
}

// Crack - runs single cracking run, based on given basePassword and rules.
// It runs password generation and schedules Pitbull instance spin up and monitoring.
func (ct *Controller) Crack(
	w http.ResponseWriter,
	r *http.Request,
) {
	var payload *models.CrackPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, err)
		return
	}

	if payload.Keyword == "" && payload.PasslistUrl == "" {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, fmt.Errorf("keyword or passlistUrl must be provided"))
		return
	}

	if payload.Keyword != "" && payload.PasslistUrl != "" {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, fmt.Errorf("only one of two arguments (keyword or passlistUrl) must be provided"))
		return
	}

	var job *models.CrackJob
	var err error

	if payload.Keyword != "" {
		job, err = ct.jobManager.CreateKeywordJob(ct.appConfig.WalletString, payload.Keyword)
		if err != nil {
			ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		job, err = ct.jobManager.CreatePasslistJob(ct.appConfig.WalletString, payload.PasslistUrl)
		if err != nil {
			ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
			return
		}
	}

	if err := ct.jobScheduler.ScheduleRun(job); err != nil {
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

// GetJobs - returns PitbullJobs based on passed filters.
func (ct *Controller) GetJobs(
	w http.ResponseWriter,
	r *http.Request,
) {
	payload := models.NewPitbullJobsListPayload()

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
