package models

import "github.com/el-mike/dogecrack/shepherd/provider"

// PitbullInstance- describes a single instance of Pitbull, backed by a provider's instance
// defined in ProviderInstance field.
type PitbullInstance struct {
	Base BaseModel `bson:",inline"`

	Name         string                  `bson:"name" json:"name"`
	Rules        []string                `bson:"rules" json:"rules"`
	Status       provider.InstanceStatus `bson:"status" json:"status"`
	ProviderName string                  `bson:"providerName" json:"providerName"`

	ProviderInstance provider.ProviderInstance `bson:"-" json:"-"`
}

func NewPitbullInstance(providerInstance *provider.ProviderInstance) *PitbullInstance {
	instance := &PitbullInstance{
		ProviderInstance: *providerInstance,
	}

	return instance
}
