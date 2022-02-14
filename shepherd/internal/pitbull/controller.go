package pitbull

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"github.com/el-mike/dogecrack/shepherd/internal/generator"
)

// Controller - responsible for handling all endpoints related to Pitbull
// instances and jobs.
type Controller struct {
	responseHelper common.ResponseHelper
	logger         *common.Logger

	appConfig *config.AppConfig

	instanceManager *InstanceManager
	jobScheduler    *Scheduler
	jobManager      *JobManager

	passwordGenerator *generator.PasswordGenerator
}

// NewController - returns new Controller instance.
func NewController() *Controller {
	instanceManager := NewInstanceManager()
	logger := common.NewLogger("Controller", os.Stdout, os.Stderr)

	return &Controller{
		responseHelper: *common.NewResponseHelper(logger),
		appConfig:      config.GetAppConfig(),

		instanceManager: instanceManager,
		jobScheduler:    NewScheduler(instanceManager),
		jobManager:      NewJobManager(instanceManager),

		passwordGenerator: generator.NewPasswordGenerator(),

		logger: logger,
	}
}

// GetActiveInstances - returns currently active instances. That includes all instances with
// status different than "FINISHED".
func (ct *Controller) GetActiveInstances(
	w http.ResponseWriter,
	r *http.Request,
) {
	instances, err := ct.instanceManager.GetActiveInstances()
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(&instances)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.responseHelper.HandleJSONResponse(w, response)
}

// GetInstance - returns an instance with given ID.
func (ct *Controller) GetInstance(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.URL.Query().Get("id")

	if id == "" {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, fmt.Errorf("Instance id was not provided"))
		return
	}

	instance, err := ct.instanceManager.GetInstanceById(id)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(instance)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.responseHelper.HandleJSONResponse(w, response)

}

// Crack - runs single cracking run, based on given basePassword and rules.
// It runs password generation and schedules Pitbull instance spin up and monitoring.
func (ct *Controller) Crack(
	w http.ResponseWriter,
	r *http.Request,
) {
	generatorResult, err := ct.passwordGenerator.Generate("testPassword1", []string{})

	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	job, err := ct.jobScheduler.ScheduleRun(generatorResult.PasslistUrl, ct.appConfig.WalletString)
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
	statusesParam := r.URL.Query().Get("statuses")

	statusesParam = strings.Trim(statusesParam, ",")

	statuses := []models.JobStatus{}

	if statusesParam != "" {
		statusesRaw := strings.Split(statusesParam, ",")

		for _, statusRaw := range statusesRaw {
			status, err := strconv.Atoi(statusRaw)
			if err != nil {
				ct.responseHelper.HandleError(w, http.StatusBadRequest, fmt.Errorf("Status: '%s' is not valid", statusRaw))
				return
			}

			statuses = append(statuses, models.JobStatus(status))
		}
	}

	jobs, err := ct.jobManager.GetJobs(statuses)
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

// RunCommand - runs a given command on Pitbull's host instance specified by
// passed ID.
func (ct *Controller) RunCommand(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.URL.Query().Get("id")

	if id == "" {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, fmt.Errorf("Instance id was not provided"))
		return
	}

	var payload struct {
		Cmd string `json:"cmd"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, err)
		return
	}

	output, err := ct.instanceManager.RunHostCommand(id, payload.Cmd)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.responseHelper.HandleJSONResponse(w, []byte(output))
}