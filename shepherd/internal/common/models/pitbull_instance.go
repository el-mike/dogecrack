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
}{
	WaitingForHost: 0,
	HostStarting:   1,
	Running:        2,
	Completed:      3,
	Failed:         4,
	Interrupted:    5,
}

var pitbullInstanceStatusByName = map[PitbullInstanceStatusEnum]string{
	PitbullInstanceStatus.WaitingForHost: "WAITING_FOR_HOST",
	PitbullInstanceStatus.HostStarting:   "HOST_STARTING",
	PitbullInstanceStatus.Running:        "RUNNING",
	PitbullInstanceStatus.Completed:      "COMPLETED",
	PitbullInstanceStatus.Failed:         "FAILED",
	PitbullInstanceStatus.Interrupted:    "INTERRUPTED",
}

// PitbullInstance - an abstract representation of Pitbull process running on some host machine.
type PitbullInstance struct {
	BaseModel `bson:",inline"`

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

// ParsePitbullStatus - parses PitbullStatus and sets PitbullInstance to Completed
// if Pitbull process is done.
func (pi *PitbullInstance) ParsePitbullStatus(status string) {
	if pi.Pitbull == nil {
		return
	}

	pi.Pitbull.ParseRawStatus(status)

	if pi.Pitbull.Status == PitbullStatus.Finished ||
		pi.Pitbull.Status == PitbullStatus.Success {
		pi.Status = PitbullInstanceStatus.Completed
	}
}

// ParsePitbullProgress - proxy method for Pitbull.ParseProgress.
func (pi *PitbullInstance) ParsePitbullProgress(progress string) error {
	return pi.Pitbull.ParseProgress(progress)
}

// Active - returns true if PitbullInstance is in one of the "active" states,
// false otherwise.
func (pi *PitbullInstance) Active() bool {
	return pi.Status == PitbullInstanceStatus.HostStarting ||
		pi.Status == PitbullInstanceStatus.Running
}

// Completed - returns true if instance's status is Completed.
func (pi *PitbullInstance) Completed() bool {
	return pi.Status == PitbullInstanceStatus.Completed
}

// AllPasswordsChecked - proxy method for Pitbull.AllPasswordsChecked.
func (pi *PitbullInstance) AllPasswordsChecked() bool {
	return pi.Pitbull.AllPasswordsChecked()
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
