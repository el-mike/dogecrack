package seeds

import (
	"github.com/el-mike/dogecrack/shepherd/internal/common/repositories"
	"os"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AppSettingsSeed - defines a seed for AppSettings.
type AppSettingsSeed struct {
	repository *repositories.AppSettingsRepository
	logger     *common.Logger
}

// NewAppSettingsSeed - returns new AppSettingsSeed instance.
func NewAppSettingsSeed() *AppSettingsSeed {
	return &AppSettingsSeed{
		repository: repositories.NewAppSettingsRepository(),
		logger:     common.NewLogger("AppSettingsSeed", os.Stdout, os.Stderr),
	}
}

// Execute - executes AppSettings seed.
func (as *AppSettingsSeed) Execute() error {
	settings := config.DefaultAppSettings
	settings.ID = primitive.NewObjectID()

	if err := as.repository.Insert(settings); err != nil {
		as.logger.Err.Println(err)

		return err
	}

	return nil
}
