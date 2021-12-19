package server

import (
	"encoding/json"
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
	passwordGenerator *generator.PasswordGenerator

	logger *log.Logger
}

// NewController - returns new Controller instance.
func NewController(manager *pitbull.PitbullManager) *Controller {
	return &Controller{
		appConfig:         config.GetAppConfig(),
		pitbullManager:    manager,
		passwordGenerator: generator.NewPasswordGenerator(),
		logger:            log.New(os.Stderr, "", 0),
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

		ct.logger.Println(err)
		return
	}

	response, err := json.Marshal(&instances)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		ct.logger.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	w.Write(response)
}

// Crack - runs single cracking run, based on given basePassword and rules.
// It runs password generation and schedules Pitbull instance spin up.
func (ct *Controller) Crack(
	w http.ResponseWriter,
	r *http.Request,
) {
	generatorResult, err := ct.passwordGenerator.Generate("testPassword1", []string{})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		ct.logger.Println(err)
		return
	}

	instance, err := ct.pitbullManager.RunInstance(generatorResult.FileUrl, ct.appConfig.WalletString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		ct.logger.Println(err)
		return
	}

	response, err := json.Marshal(instance)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		ct.logger.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	w.Write(response)
}
