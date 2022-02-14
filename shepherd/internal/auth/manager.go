package auth

import (
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"github.com/el-mike/dogecrack/shepherd/internal/users"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Manager - responsible for managing all auth-related operations.
type Manager struct {
	usersManager *users.Manager
	cache        *Cache
}

// NewManager - returns new Manager instance.
func NewManager() *Manager {
	appConfig := config.GetAppConfig()

	return &Manager{
		usersManager: users.NewManager(),
		cache:        NewCache(appConfig.SessionExpiration),
	}
}

// GetMe - returns currently logged-in User.
func (mg *Manager) GetUser(userId string) (*models.User, error) {
	return mg.usersManager.GetUserById(userId)
}

// Login - logs in given user and saves sessionId to cache.
func (mg *Manager) Login(username, password string) (string, error) {
	user, err := mg.usersManager.GetUserByName(username)
	if err != nil {
		return "", NewUserNotFoundError(username)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", NewIncorrectPasswordError()
	}

	userId := user.ID.Hex()
	sessionId := uuid.New().String()

	if err := mg.cache.SaveSessionId(sessionId, userId); err != nil {
		return "", nil
	}

	return sessionId, nil
}

// Logout - logs out given user.
func (mg *Manager) Logout(sessionId string) error {
	return mg.cache.DeleteSessionId(sessionId)
}
