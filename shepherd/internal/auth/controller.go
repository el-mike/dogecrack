package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/config"
)

// Controller - responsible for auth-related routes.
type Controller struct {
	responseHelper *common.ResponseHelper
	logger         *common.Logger

	manager *Manager

	expiration time.Duration
}

// NewController - returns a new Controller instance.
func NewController() *Controller {
	config := config.GetAppConfig()
	logger := common.NewLogger("AuthController", os.Stdout, os.Stderr)

	return &Controller{
		responseHelper: common.NewResponseHelper(logger),
		logger:         logger,

		manager: NewManager(),

		expiration: config.SessionExpiration,
	}
}

// Login - performs login operation and sets proper cookie for the client.
func (ct *Controller) Login(
	w http.ResponseWriter,
	r *http.Request,
) {
	var creds *models.UserCredentials

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, err)
		return
	}

	if creds.Name == "" || creds.Password == "" {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, fmt.Errorf("Username or password not provided"))
		return
	}

	sessionId, err := ct.manager.Login(creds.Name, creds.Password)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionIdCookie,
		Value:    sessionId,
		Path:     "/",
		Expires:  time.Now().Add(ct.expiration),
		HttpOnly: true,
	})

	ct.responseHelper.HandleEmptyResponse(w)
}
