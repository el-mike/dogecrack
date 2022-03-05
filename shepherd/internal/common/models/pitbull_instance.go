package models

import (
	"github.com/el-mike/dogecrack/shepherd/internal/common/host"
	"github.com/el-mike/dogecrack/shepherd/internal/vast/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PitbullInstanceStatusEnum - enum describing the status of PitbullInstance.
type PitbullInstanceStatusEnum int8

// PitbullInstanceStatus - "namespaced" struct containing PitbullInstanceStatus values.
var PitbullInstanceStatus = struct {
	WaitingForHost PitbullInstanceStatusEnum
	HostStarting   PitbullInstanceStatusEnum
	Running        PitbullInstanceStatusEnum
	Completed      PitbullInstanceStatusEnum
	Failed         PitbullInstanceStatusEnum
	Interrupted    PitbullInstanceStatusEnum
	Success        PitbullInstanceStatusEnum
}{
	WaitingForHost: 0,
	HostStarting:   1,
	Running:        2,
	Completed:      3,
	Failed:         4,
	Interrupted:    5,
	Success:        6,
}

var pitbullInstanceStatusByName = map[PitbullInstanceStatusEnum]string{
	PitbullInstanceStatus.WaitingForHost: "WAITING_FOR_HOST",
	PitbullInstanceStatus.HostStarting:   "HOST_STARTING",
	PitbullInstanceStatus.Running:        "RUNNING",
	PitbullInstanceStatus.Completed:      "COMPLETED",
	PitbullInstanceStatus.Failed:         "FAILED",
	PitbullInstanceStatus.Interrupted:    "INTERRUPTED",
	PitbullInstanceStatus.Success:        "SUCCESS",
}

// PitbullInstance - an abstract representation of Pitbull process running on some host machine.
type PitbullInstance struct {
	BaseModel `bson:",inline"`

	JobID primitive.ObjectID `bson:"jobId" json:"jobId"`

	WalletString string `bson:"walletString" json:"walletString"`
	PasslistUrl  string `bson:"passlistUrl" json:"passlistUrl"`

	StartedAt   NullableTime `bson:"startedAt" json:"startedAt"`
	CompletedAt NullableTime `bson:"completedAt" json:"completedAt"`

	Status PitbullInstanceStatusEnum `bson:"status" json:"status"`

	Pitbull *Pitbull `bson:"pitbull" json:"pitbull"`

	ProviderName    string            `bson:"providerName" json:"providerName"`
	HostInstance    host.HostInstance `bson:"-" json:"hostInstance"`
	HostInstanceRaw bson.Raw          `bson:"hostInstanceRaw" json:"-"`
}

type marshalablePitbullInstance PitbullInstance

// NewPitbullInstance - returns new PitbullInstance instance.
func NewPitbullInstance(host host.HostInstance, passlistUrl, walletString string) *PitbullInstance {
	instance := &PitbullInstance{
		PasslistUrl:  passlistUrl,
		WalletString: walletString,
		Status:       PitbullInstanceStatus.WaitingForHost,
		HostInstance: host,
		Pitbull:      NewPitbull(),
	}

	instance.ID = primitive.NewObjectID()

	return instance
}

// SetHost - attaches HostInstance to given PitbullInstance.
func (pi *PitbullInstance) SetHost(hostInstance host.HostInstance) {
	pi.HostInstance = hostInstance
	pi.ProviderName = hostInstance.ProviderName()
}

// Active - returns true if PitbullInstance is in one of the "active" states,
// false otherwise.
func (pi *PitbullInstance) Active() bool {
	return pi.Status == PitbullInstanceStatus.HostStarting ||
		pi.Status == PitbullInstanceStatus.Running
}

// UnmarshalBSON - Unmarshaler interface implementation.
func (pi *PitbullInstance) UnmarshalBSON(data []byte) error {
	if err := bson.Unmarshal(data, (*marshalablePitbullInstance)(pi)); err != nil {
		return err
	}

	if pi.ProviderName == models.ProviderName {
		vastInstance := &models.Instance{}
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
