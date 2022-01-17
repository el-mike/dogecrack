package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/el-mike/dogecrack/shepherd/config"
	"github.com/el-mike/dogecrack/shepherd/generator"
	"github.com/el-mike/dogecrack/shepherd/pitbull"
)

type ControllerFn func(w http.ResponseWriter, r *http.Request)

// Controller - instance responsible for handling API requests.
type Controller struct {
	appConfig *config.AppConfig

	pitbullManager    *pitbull.PitbullManager
	pitbullMonitor    *pitbull.PitbullMonitor
	passwordGenerator *generator.PasswordGenerator

	errorLogger *log.Logger
}

// NewController - returns new Controller instance.
func NewController(manager *pitbull.PitbullManager) *Controller {
	return &Controller{
		appConfig: config.GetAppConfig(),

		pitbullManager:    manager,
		pitbullMonitor:    pitbull.NewPitbullMonitor(manager),
		passwordGenerator: generator.NewPasswordGenerator(),

		errorLogger: log.New(os.Stderr, "[Controller][Error]: ", log.Ldate|log.Ltime),
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

		ct.errorLogger.Println(err)
		return
	}

	response, err := json.Marshal(&instances)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		ct.errorLogger.Println(err)
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
// It runs password generation and schedules Pitbull instance spin up.
func (ct *Controller) Crack(
	w http.ResponseWriter,
	r *http.Request,
) {
	generatorResult, err := ct.passwordGenerator.Generate("testPassword1", []string{})

	if err != nil {
		ct.handleError(w, http.StatusInternalServerError, err)
		return
	}

	instance, err := ct.pitbullManager.RunInstance(generatorResult.FileUrl, ct.appConfig.WalletString)
	if err != nil {
		ct.handleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.pitbullMonitor.RunMonitoring(instance.ID.Hex())

	response, err := json.Marshal(instance)
	if err != nil {
		ct.handleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.handleJSONResponse(w, response)
}

// handleInternalError - helper function for returning error as JSON.
func (ct *Controller) handleError(w http.ResponseWriter, status int, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	ct.errorLogger.Println(err)

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
