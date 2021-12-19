package models

import (
	"github.com/el-mike/dogecrack/shepherd/provider"
	"github.com/el-mike/dogecrack/shepherd/vast"
	"go.mongodb.org/mongo-driver/bson"
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

	ProviderInstance    provider.ProviderInstance `bson:"-" json:"providerInstance"`
	ProviderInstanceRaw bson.Raw                  `bson:"providerInstanceRaw" json:"-"`
}

type marshalablePitbullInstance PitbullInstance

// NewPitbullInstance - returns new PitbullInstance instance.
func NewPitbullInstance(providerInstance provider.ProviderInstance) *PitbullInstance {
	instance := &PitbullInstance{
		ProviderInstance: providerInstance,
		ProviderName:     providerInstance.ProviderName(),
	}

	instance.ID = primitive.NewObjectID()

	return instance
}

func (pi *PitbullInstance) UnmarshalBSON(data []byte) error {
	if err := bson.Unmarshal(data, (*marshalablePitbullInstance)(pi)); err != nil {
		return err
	}

	if pi.ProviderName == vast.ProviderName {
		vastInstance := &vast.VastInstance{}
		if err := bson.Unmarshal(pi.ProviderInstanceRaw, &vastInstance); err != nil {
			return err
		}

		pi.ProviderInstance = vastInstance
	}

	return nil
}

func (pi *PitbullInstance) MarshalBSON() ([]byte, error) {
	providerRaw, err := bson.Marshal(pi.ProviderInstance)
	if err != nil {
		return nil, err
	}

	pi.ProviderInstanceRaw = providerRaw

	return bson.Marshal((*marshalablePitbullInstance)(pi))
}
