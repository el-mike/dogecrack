package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/common/api"
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/config"
)

// Controller - responsible for auth-related routes.
type Controller struct {
	responseHelper *api.ResponseHelper
	logger         *common.Logger

	manager *Manager

	expiration time.Duration
}

// NewController - returns a new Controller instance.
func NewController() *Controller {
	config := config.GetAppConfig()
	logger := common.NewLogger("AuthController", os.Stdout, os.Stderr)

	return &Controller{
		responseHelper: api.NewResponseHelper(logger),
		logger:         logger,

		manager: NewManager(),

		expiration: config.SessionExpiration,
	}
}

// Me - returns a User related to current session sent via cookie.
func (ct *Controller) Me(
	w http.ResponseWriter,
	r *http.Request,
) {
	contextUser := r.Context().Value(contextUserKey).(*ContextUser)

	user, err := ct.manager.GetUser(contextUser.userId)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	userResponse := models.NewUserResponse(user)

	response, err := json.Marshal(&userResponse)
	if err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	ct.responseHelper.HandleJSONResponse(w, response)
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

	if creds.Username == "" || creds.Password == "" {
		ct.responseHelper.HandleError(w, http.StatusBadRequest, fmt.Errorf("Username or password not provided"))
		return
	}

	sessionId, err := ct.manager.Login(creds.Username, creds.Password)
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

// Logout - logs user out by deleting the cache entry and cookie.
func (ct *Controller) Logout(
	w http.ResponseWriter,
	r *http.Request,
) {
	contextUser := r.Context().Value(contextUserKey).(*ContextUser)

	if err := ct.manager.Logout(contextUser.sessionId); err != nil {
		ct.responseHelper.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionIdCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	ct.responseHelper.HandleEmptyResponse(w)
}
