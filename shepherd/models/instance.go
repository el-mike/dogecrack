package models

import (
	"github.com/el-mike/dogecrack/shepherd/provider"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PitbullInstance- describes a single instance of Pitbull, backed by a provider's instance
// defined in ProviderInstance field.
type PitbullInstance struct {
	BaseModel `bson:",inline"`

	Name         string                  `bson:"name" json:"name"`
	Rules        []string                `bson:"rules" json:"rules"`
	Status       provider.InstanceStatus `bson:"status" json:"status"`
	ProviderName string                  `bson:"providerName" json:"providerName"`

	ProviderInstance provider.ProviderInstance `bson:"-" json:"-"`
}

func NewPitbullInstance(providerInstance provider.ProviderInstance) *PitbullInstance {
	instance := &PitbullInstance{
		ProviderInstance: providerInstance,
		ProviderName:     providerInstance.ProviderName(),
	}

	instance.ID = primitive.NewObjectID()

	return instance
}
