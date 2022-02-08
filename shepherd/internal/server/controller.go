package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"github.com/el-mike/dogecrack/shepherd/internal/generator"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/models"
)

type ControllerFn func(w http.ResponseWriter, r *http.Request)

// Controller - instance responsible for handling API requests.
type Controller struct {
	appConfig *config.AppConfig

	pitbullManager   *pitbull.Manager
	pitbullScheduler *pitbull.Scheduler
	jobService       *pitbull.JobService

	passwordGenerator *generator.PasswordGenerator

	logger *common.Logger
}

// NewController - returns new Controller instance.
func NewController(manager *pitbull.Manager) *Controller {
	return &Controller{
		appConfig: config.GetAppConfig(),

		pitbullManager:   manager,
		pitbullScheduler: pitbull.NewScheduler(),
		jobService:       pitbull.NewJobService(),

		passwordGenerator: generator.NewPasswordGenerator(),

		logger: common.NewLogger("Controller", os.Stdout, os.Stderr),
	}
}

// GetHealth - handles health check.
func (ct *Controller) GetHealth(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(http.StatusOK)
}

// GetActiveInstances - returns currently active instances. That includes all instances with
// status different than "FINISHED".
func (ct *Controller) GetActiveInstances(
	w http.ResponseWriter,
	r *http.Request,
) {
	instances, err := ct.pitbullManager.GetActiveInstances()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		ct.logger.Err.Println(err)
		return
	}

	response, err := json.Marshal(&instances)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		ct.logger.Err.Println(err)
		return
	}

	ct.handleJSONResponse(w, response)
}

// GetInstance - returns an instance with given ID.
func (ct *Controller) GetInstance(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.URL.Query().Get("id")

	if id == "" {
		ct.handleError(w, http.StatusBadRequest, fmt.Errorf("Instance id was not provided"))
		return
	}

	instance, err := ct.pitbullManager.GetInstanceById(id)
	if err != nil {
		ct.handleError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(instance)
	if err != nil {
		ct.handleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.handleJSONResponse(w, response)

}

// Crack - runs single cracking run, based on given basePassword and rules.
// It runs password generation and schedules Pitbull instance spin up and monitoring.
func (ct *Controller) Crack(
	w http.ResponseWriter,
	r *http.Request,
) {
	generatorResult, err := ct.passwordGenerator.Generate("testPassword1", []string{})

	if err != nil {
		ct.handleError(w, http.StatusInternalServerError, err)
		return
	}

	instance, err := ct.pitbullManager.CreateInstance(generatorResult.PasslistUrl, ct.appConfig.WalletString)
	if err != nil {
		ct.handleError(w, http.StatusInternalServerError, err)
		return
	}

	job, err := ct.pitbullScheduler.ScheduleRun(instance)
	if err != nil {
		ct.handleError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(job)
	if err != nil {
		ct.handleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.handleJSONResponse(w, response)
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
				ct.handleError(w, http.StatusBadRequest, fmt.Errorf("Status: '%s' is not valid", statusRaw))
				return
			}

			statuses = append(statuses, models.JobStatus(status))
		}
	}

	jobs, err := ct.jobService.GetJobs(statuses)
	if err != nil {
		ct.handleError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(jobs)
	if err != nil {
		ct.handleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.handleJSONResponse(w, response)
}

// RunCommand - runs a given command on Pitbull's host instance specified by
// passed ID.
func (ct *Controller) RunCommand(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.URL.Query().Get("id")

	if id == "" {
		ct.handleError(w, http.StatusBadRequest, fmt.Errorf("Instance id was not provided"))
		return
	}

	var payload struct {
		Cmd string `json:"cmd"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		ct.handleError(w, http.StatusBadRequest, err)
		return
	}

	output, err := ct.pitbullManager.RunHostCommand(id, payload.Cmd)
	if err != nil {
		ct.handleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.handleJSONResponse(w, []byte(output))
}

// handleInternalError - helper function for returning error as JSON.
func (ct *Controller) handleError(w http.ResponseWriter, status int, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	common.WithStackTrace(ct.logger.Err, err)

	apiError := NewApiError(status, err)

	errorResponse, err := json.Marshal(apiError)
	if err != nil {
		// Do nothing - status 500 will be returned anyway.
	}

	w.Write(errorResponse)
}

// handleJSONResponse - helper function for returning success response as JSON.
func (ct *Controller) handleJSONResponse(w http.ResponseWriter, response []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(response)
}
