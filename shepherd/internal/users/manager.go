package users

import (
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/users/repositories"
)

// Manager - responsible for all User-related logic.
type Manager struct {
	repository *repositories.UserRepository
}

// NewManager - returns new Manager instance.
func NewManager() *Manager {
	return &Manager{
		repository: repositories.NewUserRepository(),
	}
}

// SaveUser - inserts a single User to the DB.
func (mg *Manager) SaveUser(user *models.User) error {
	return mg.repository.Insert(user)
}

// GetUserByName - returns a single User with given name.
func (mg *Manager) GetUserByName(name string) (*models.User, error) {
	return mg.repository.GetByName(name)
}