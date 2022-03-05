package core

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/common/api"
)

// Controller - instance responsible for handling general, app-related endpoints.
type Controller struct {
	generalService *GeneralService

	responseHelper api.ResponseHelper
}

// NewController - returns new Controller instance.
func NewController() *Controller {
	logger := common.NewLogger("AppController", os.Stdout, os.Stderr)

	return &Controller{
		generalService: NewGeneralService(),

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
	enums := ct.generalService.GetEnums()

	response, err := json.Marshal(&enums)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.responseHelper.HandleJSONResponse(w, response)
}

// GetStatistics - returns application's statistics.
func (ct *Controller) GetStatistics(
	w http.ResponseWriter,
	r *http.Request,
) {
	statistics, err := ct.generalService.GetStatistics()
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(&statistics)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.responseHelper.HandleJSONResponse(w, response)
}
