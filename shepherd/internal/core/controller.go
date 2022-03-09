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

// GetSettings - returns application's settings.
func (ct *Controller) GetSettings(
	w http.ResponseWriter,
	r *http.Request,
) {
	settings, err := ct.generalService.GetAppSettings()
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(&settings)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.responseHelper.HandleJSONResponse(w, response)
}

// UpdateSettings - updates current app settings.
func (ct *Controller) UpdateSettings(
	w http.ResponseWriter,
	r *http.Request,
) {
	var payload *models.AppSettingsPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, err)
		return
	}

	settings, err := ct.generalService.UpdateAppSettings(payload)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(&settings)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.responseHelper.HandleJSONResponse(w, response)
}
