package api

import (
	"encoding/json"
	"net/http"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
)

// ResponseHelper - a helper entity, containg a common Controller logic
// for handling requests and responses.
type ResponseHelper struct {
	logger *common.Logger
}

// NewResponseHelper - return new BaseController instance.
func NewResponseHelper(logger *common.Logger) *ResponseHelper {
	return &ResponseHelper{
		logger: logger,
	}
}

// handleInternalError - helper function for returning error as JSON.
func (ct *ResponseHelper) HandleError(w http.ResponseWriter, status int, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	common.WithStackTrace(ct.logger.Err, err)

	apiError := NewApiError(status, err)

	errorResponse, err := json.Marshal(apiError)
	if err != nil {
		// Do nothing - status 500 will be returned anyway.
	}

	w.Write(errorResponse)
}

// handleJSONResponse - helper function for returning success response as JSON.
func (ct *ResponseHelper) HandleJSONResponse(w http.ResponseWriter, response []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(response)
}

// HandleEmptyResponse - helper function for returning 200 OK status.
func (ct *ResponseHelper) HandleEmptyResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
