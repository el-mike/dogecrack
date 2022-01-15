package models

import (
	"errors"
	"strconv"
	"strings"

	"github.com/el-mike/dogecrack/shepherd/host"
	"github.com/el-mike/dogecrack/shepherd/vast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PitbullStatus - enum describing the status of a Pitbull process.
type PitbullStatus int8

const (
	Waiting PitbullStatus = iota
	Running
	Finished
	Success
)

// ProgressInfo - helper struct describing Pitbull progress.
type ProgressInfo struct {
	Checked int64 `bson:"checked" json:"checked"`
	Total   int64 `bson:"total" json:"total"`
}

// PitbullInstance- describes a single instance of PitbullInstance, backed by a host instance
// defined in ProviderInstance field.
type PitbullInstance struct {
	BaseModel `bson:",inline"`

	Name     string        `bson:"name" json:"name"`
	Rules    []string      `bson:"rules" json:"rules"`
	Status   PitbullStatus `bson:"status" json:"status"`
	Progress *ProgressInfo `bson:"progress" json:"progress"`

	ProviderName string `bson:"providerName" json:"providerName"`

	HostInstance    host.HostInstance `bson:"-" json:"hostInstance"`
	HostInstanceRaw bson.Raw          `bson:"hostInstanceRaw" json:"-"`
}

type marshalablePitbullInstance PitbullInstance

// NewPitbullInstance - returns new PitbullInstance instance.
func NewPitbullInstance(hostInstance host.HostInstance) *PitbullInstance {
	instance := &PitbullInstance{
		Status:       Waiting,
		HostInstance: hostInstance,
		ProviderName: hostInstance.ProviderName(),
	}

	instance.ID = primitive.NewObjectID()

	return instance
}

// Completed - returns true if all passwords have been checked for given Pitbull instance,
// false otherwise.
func (pi *PitbullInstance) Completed() bool {
	if pi.Progress == nil || pi.Progress.Total == 0 {
		return false
	}

	return pi.Progress.Checked == pi.Progress.Total
}

// SetStatus - helper function for converting raw status command output into
// PitbullStatus enum.
func (pi *PitbullInstance) SetStatus(rawStatus string) {
	if strings.Contains(rawStatus, "SUCCESS") {
		pi.Status = Success
	} else if strings.Contains(rawStatus, "RUNNING") {
		pi.Status = Running
	} else if strings.Contains(rawStatus, "FINISHED") {
		pi.Status = Finished
	} else {
		pi.Status = Waiting
	}
}

// SetProgress - helper function for converting raw progress command output into
// ProgressInfo struct.
func (pi *PitbullInstance) SetProgress(rawProgress string) error {
	if pi.Progress == nil {
		pi.Progress = &ProgressInfo{}
	}
	parts := strings.Split(rawProgress, "of")

	if parts == nil || parts[0] == "" || parts[1] == "" {
		return errors.New("Progress malformed: " + rawProgress)
	}

	checkedStr := strings.TrimSpace(parts[0])
	totalStr := strings.TrimSpace(parts[1])

	checked, err := strconv.Atoi(checkedStr)
	if err != nil {
		return err
	}

	total, err := strconv.Atoi(totalStr)
	if err != nil {
		return err
	}

	pi.Progress.Checked = int64(checked)
	pi.Progress.Total = int64(total)

	return nil
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