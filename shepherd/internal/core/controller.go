package core

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/common/api"
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
)

// Controller - instance responsible for handling general, app-related endpoints.
type Controller struct {
	responseHelper api.ResponseHelper
}

// NewController - returns new Controller instance.
func NewController() *Controller {
	logger := common.NewLogger("AppController", os.Stdout, os.Stderr)

	return &Controller{
		responseHelper: *api.NewResponseHelper(logger),
	}
}

// GetHealth - handles health check.
func (ct *Controller) GetHealth(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(http.StatusOK)
}

// GetEnums - returns application enums configuration.
func (ct *Controller) GetEnums(
	w http.ResponseWriter,
	r *http.Request,
) {
	enums := models.GetAppEnums()

	response, err := json.Marshal(&enums)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.responseHelper.HandleJSONResponse(w, response)
}
