package models

import (
	"github.com/el-mike/dogecrack/shepherd/host"
	"github.com/el-mike/dogecrack/shepherd/vast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PitbullStatus - enum describing the status of a Pitbull process.
type PitbullStatus int8

const (
	Running PitbullStatus = iota
	Finished
	Success
)

// PitbullInstance- describes a single instance of PitbullInstance, backed by a host instance
// defined in ProviderInstance field.
type PitbullInstance struct {
	BaseModel `bson:",inline"`

	Name         string        `bson:"name" json:"name"`
	Rules        []string      `bson:"rules" json:"rules"`
	Status       PitbullStatus `bson:"status" json:"status"`
	ProviderName string        `bson:"providerName" json:"providerName"`

	HostInstance    host.HostInstance `bson:"-" json:"hostInstance"`
	HostInstanceRaw bson.Raw          `bson:"hostInstanceRaw" json:"-"`
}

type marshalablePitbullInstance PitbullInstance

// NewPitbullInstance - returns new PitbullInstance instance.
func NewPitbullInstance(hostInstance host.HostInstance) *PitbullInstance {
	instance := &PitbullInstance{
		HostInstance: hostInstance,
		ProviderName: hostInstance.ProviderName(),
	}

	instance.ID = primitive.NewObjectID()

	return instance
}

// UnmarshalBSON - Unmarshaler interface implementation.
func (pi *PitbullInstance) UnmarshalBSON(data []byte) error {
	if err := bson.Unmarshal(data, (*marshalablePitbullInstance)(pi)); err != nil {
		return err
	}

	if pi.ProviderName == vast.ProviderName {
		vastInstance := &vast.VastInstance{}
		if err := bson.Unmarshal(pi.HostInstanceRaw, &vastInstance); err != nil {
			return err
		}

		pi.HostInstance = vastInstance
	}

	return nil
}

// MarshalBSON - Marshaler interface implementation.
func (pi *PitbullInstance) MarshalBSON() ([]byte, error) {
	providerRaw, err := bson.Marshal(pi.HostInstance)
	if err != nil {
		return nil, err
	}

	pi.HostInstanceRaw = providerRaw

	return bson.Marshal((*marshalablePitbullInstance)(pi))
}
