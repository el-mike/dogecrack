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
	"github.com/el-mike/dogecrack/shepherd/internal/generator"
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

	passwordGenerator *generator.PasswordGenerator
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

		passwordGenerator: generator.NewPasswordGenerator(),

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

	if payload.Keyword == "" {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, fmt.Errorf("Keyword is required"))
		return
	}

	generatorResult, err := ct.passwordGenerator.Generate(payload.Keyword, payload.Rules)

	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	job, err := ct.jobScheduler.ScheduleRun(generatorResult.Keyword, generatorResult.PasslistUrl, ct.appConfig.WalletString)
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
