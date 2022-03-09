package seeds

import (
	"os"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/users"
	"golang.org/x/crypto/bcrypt"
)

// UsersSeed - defines a seed for Users collection.
type UsersSeed struct {
	repository *users.Repository
	logger     *common.Logger
}

// NewUsersSeed - returns new UsersSeed instance.
func NewUsersSeed() *UsersSeed {
	return &UsersSeed{
		repository: users.NewRepository(),
		logger:     common.NewLogger("UsersSeed", os.Stdout, os.Stderr),
	}
}

// Execute - executes the seed.
func (us *UsersSeed) Execute() error {
	name := os.Getenv("ADMIN_NAME")
	password := os.Getenv("ADMIN_PASSWORD")

	existingAdmin, err := us.repository.GetByName(name)
	if err != nil {
		return err
	}

	// If admin already exists, do nothing.
	if existingAdmin != nil {
		us.logger.Err.Println("Admin already exists")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	admin := models.NewUser(name, string(hashedPassword))
	admin.Role = "ADMIN"

	return us.repository.Insert(admin)
}
