package core

import (
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/crack"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull"
)

// GeneralService - responsible for general, often app-wide operations.
type GeneralService struct {
	crackJobRepository        *crack.JobRepository
	pitbullInstanceRepository *pitbull.InstanceRepository
	appSettingsRepository     *AppSettingsRepository
}

// NewGeneralService - returns new GeneralService instance.
func NewGeneralService() *GeneralService {
	return &GeneralService{
		crackJobRepository:        crack.NewJobRepository(),
		pitbullInstanceRepository: pitbull.NewInstanceRepository(),
		appSettingsRepository:     NewAppSettingsRepository(),
	}
}

// GetEnums - returns app's reusable enums.
func (gs *GeneralService) GetEnums() *models.AppEnums {
	return models.GetAppEnums()
}

// GetStatistics - returns app's overall statistics.
func (gs *GeneralService) GetStatistics() (*models.AppStatistics, error) {
	jobsStatistics, err := gs.crackJobRepository.GetStatistics()
	if err != nil {
		return nil, err
	}

	instancesStatistics, err := gs.pitbullInstanceRepository.GetStatistics()
	if err != nil {
		return nil, err
	}

	statistics := &models.AppStatistics{
		CrackJobs:        jobsStatistics,
		PitbullInstances: instancesStatistics,
	}

	return statistics, nil
}

// GetAppSettings - returns app' editable settings.
func (gs *GeneralService) GetAppSettings() (*models.AppSettings, error) {
	return gs.appSettingsRepository.GetAppSettings()
}

// UpdateAppSettings - updates current AppSettings.
func (gs *GeneralService) UpdateAppSettings(payload *models.AppSettingsPayload) (*models.AppSettings, error) {
	settings, err := gs.appSettingsRepository.GetAppSettings()
	if err != nil {
		return nil, err
	}

	settings.Merge(payload)

	if err := gs.appSettingsRepository.Update(settings); err != nil {
		return nil, err
	}

	return settings, nil
}
