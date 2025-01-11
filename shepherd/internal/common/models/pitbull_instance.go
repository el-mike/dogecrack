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
	Success        PitbullInstanceStatusEnum
}{
	WaitingForHost: 0,
	HostStarting:   1,
	Running:        2,
	Completed:      3,
	Failed:         4,
	Success:        5,
}

var pitbullInstanceStatusByName = map[PitbullInstanceStatusEnum]string{
	PitbullInstanceStatus.WaitingForHost: "WAITING_FOR_HOST",
	PitbullInstanceStatus.HostStarting:   "HOST_STARTING",
	PitbullInstanceStatus.Running:        "RUNNING",
	PitbullInstanceStatus.Completed:      "COMPLETED",
	PitbullInstanceStatus.Failed:         "FAILED",
	PitbullInstanceStatus.Success:        "SUCCESS",
}

type PitbullRunPayload struct {
	WalletString string `bson:"walletString" json:"walletString"`
	Tokenlist    string `bson:"tokenlist" json:"tokenlist"`
	PasslistUrl  string `bson:"passlistUrl" json:"passlistUrl"`
	SkipCount    int64  `bson:"skipCount" json:"skipCount"`
}

// PitbullInstance - an abstract representation of Pitbull process running on some host machine.
type PitbullInstance struct {
	BaseModel `bson:",inline"`

	JobId      primitive.ObjectID `bson:"jobId" json:"jobId"`
	RunPayload *PitbullRunPayload `bson:"runPayload" json:"runPayload"`

	StartedAt   NullableTime `bson:"startedAt" json:"startedAt"`
	CompletedAt NullableTime `bson:"completedAt" json:"completedAt"`

	Status PitbullInstanceStatusEnum `bson:"status" json:"status"`

	Pitbull *Pitbull `bson:"pitbull" json:"pitbull"`

	ProviderName    string            `bson:"providerName" json:"providerName"`
	HostInstance    host.HostInstance `bson:"-" json:"hostInstance"`
	HostInstanceRaw bson.Raw          `bson:"hostInstanceRaw" json:"-"`

	FailReason string `bson:"failReason" json:"failReason"`
}

type marshalablePitbullInstance PitbullInstance

// NewPitbullInstance - returns new PitbullInstance instance.
func NewPitbullInstance(host host.HostInstance, runPayload *PitbullRunPayload) *PitbullInstance {
	instance := &PitbullInstance{
		Status:       PitbullInstanceStatus.WaitingForHost,
		HostInstance: host,
		RunPayload:   runPayload,
		Pitbull:      NewPitbull(runPayload.SkipCount),
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
