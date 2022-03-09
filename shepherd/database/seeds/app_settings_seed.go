package seeds

import (
	"os"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"github.com/el-mike/dogecrack/shepherd/internal/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AppSettingsSeed - defines a seed for AppSettings.
type AppSettingsSeed struct {
	repository *core.AppSettingsRepository
	logger     *common.Logger
}

// NewAppSettingsSeed - returns new AppSettingsSeed instance.
func NewAppSettingsSeed() *AppSettingsSeed {
	return &AppSettingsSeed{
		repository: core.NewAppSettingsRepository(),
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
