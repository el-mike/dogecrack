package pitbull

import (
	"github.com/el-mike/dogecrack/shepherd/internal/common/repositories"
)

const (
	DefaultOfferFilter = "dph<3"
)

// VastSearchCriteriaProvider - provides search criteria from app settings.
type VastSearchCriteriaProvider struct {
	appSettingsRepository *repositories.AppSettingsRepository
}

// NewVastSearchCriteriaProvider - returns new VastSearchCriteriaProvider instance
func NewVastSearchCriteriaProvider() *VastSearchCriteriaProvider {
	return &VastSearchCriteriaProvider{
		appSettingsRepository: repositories.NewAppSettingsRepository(),
	}
}

func (vp *VastSearchCriteriaProvider) GetSearchCriteria() string {
	settings, err := vp.appSettingsRepository.GetAppSettings()
	if err != nil || settings.VastInstanceSearchCriteria == "" {
		return DefaultOfferFilter
	}

	return settings.VastInstanceSearchCriteria
}
