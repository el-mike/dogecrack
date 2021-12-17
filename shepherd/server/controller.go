package server

import (
	"net/http"

	"github.com/el-mike/dogecrack/shepherd/pitbull"
)

type ControllerFn func(w http.ResponseWriter, r *http.Request)

// Controller - instance responsible for handling API requests.
type Controller struct {
	pitbullManager pitbull.PitbullManager
	pitbullClient  pitbull.PitbullClient
}

func NewController(manager pitbull.PitbullManager, client pitbull.PitbullClient) *Controller {
	return &Controller{
		pitbullManager: manager,
		pitbullClient:  client,
	}
}

// GetHealth - handles health check.
func (ct *Controller) GetHealth(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(http.StatusOK)
}
