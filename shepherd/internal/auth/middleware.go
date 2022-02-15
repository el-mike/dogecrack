package auth

import (
	"context"
	"net/http"
	"os"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/config"
)

// Middleware - responsible for auth validation.
type Middleware struct {
	responseHelper *common.ResponseHelper
	logger         *common.Logger
	cache          *Cache
}

// NewMiddleware - returns new Middleware instance.
func NewMiddleware() *Middleware {
	logger := common.NewLogger("AuthMiddleware", os.Stdout, os.Stderr)

	appConfig := config.GetAppConfig()

	return &Middleware{
		responseHelper: common.NewResponseHelper(logger),
		logger:         logger,
		cache:          NewCache(appConfig.SessionExpiration),
	}
}

// Middleware - gorilla/mux middleware interface implementation.
func (md *Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionIdCookie)
		if err != nil {
			if err == http.ErrNoCookie {
				md.responseHelper.HandleError(w, http.StatusUnauthorized, NewUserUnathorizedError())
				return
			}

			md.responseHelper.HandleError(w, http.StatusBadRequest, err)
			return
		}

		sessionId := cookie.Value

		userId, err := md.cache.GetUserBySessionId(sessionId)
		if err != nil {
			md.responseHelper.HandleError(w, http.StatusInternalServerError, err)
			return
		}

		contextUser := &ContextUser{
			sessionId: sessionId,
			userId:    userId,
		}

		ctx := context.WithValue(r.Context(), contextUserKey, contextUser)

		if userId == "" {
			md.responseHelper.HandleError(w, http.StatusUnauthorized, NewSessionExpiredError())
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
