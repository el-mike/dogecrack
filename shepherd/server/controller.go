package server

import (
	"fmt"
	"net/http"

	"github.com/el-mike/dogecrack/shepherd/config"
	"github.com/el-mike/dogecrack/shepherd/generator"
	"github.com/el-mike/dogecrack/shepherd/pitbull"
)

type ControllerFn func(w http.ResponseWriter, r *http.Request)

// Controller - instance responsible for handling API requests.
type Controller struct {
	appConfig      *config.AppConfig
	pitbullManager pitbull.PitbullManager
	pitbullClient  pitbull.PitbullClient

	passwordGenerator *generator.PasswordGenerator
}

// NewController - returns new Controller instance.
func NewController(manager pitbull.PitbullManager, client pitbull.PitbullClient) *Controller {
	return &Controller{
		appConfig:         config.GetAppConfig(),
		pitbullManager:    manager,
		pitbullClient:     client,
		passwordGenerator: generator.NewPasswordGenerator(),
	}
}

// GetHealth - handles health check.
func (ct *Controller) GetHealth(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(http.StatusOK)
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
	}

	instance, err := ct.pitbullManager.RunInstance(generatorResult.FileUrl, ct.appConfig.WalletString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Print(instance)
}
