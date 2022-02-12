package server

import (
	"net/http"
	"os"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
)

// Controller - instance responsible for handling general, app-related endpoints.
type Controller struct {
	responseHelper common.ResponseHelper
}

// NewController - returns new Controller instance.
func NewController() *Controller {
	logger := common.NewLogger("AppController", os.Stdout, os.Stderr)

	return &Controller{
		responseHelper: *common.NewResponseHelper(logger),
	}
}

// GetHealth - handles health check.
func (ct *Controller) GetHealth(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(http.StatusOK)
}
