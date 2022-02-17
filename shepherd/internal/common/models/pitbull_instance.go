package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/el-mike/dogecrack/shepherd/internal/common/host"
	"github.com/el-mike/dogecrack/shepherd/internal/vast/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PitbullStatus - enum describing the status of a Pitbull process.
type PitbullStatus int8

const (
	WaitingForHost PitbullStatus = iota
	HostStarting
	Waiting
	Running
	Finished
	Success
	Interrupted
	Failed
)

var pitbullStatusNames = map[PitbullStatus]string{
	WaitingForHost: "WAITING_FOR_HOST",
	HostStarting:   "HOST_STARTING",
	Waiting:        "WAITING",
	Running:        "RUNNING",
	Finished:       "FINISHED",
	Success:        "SUCCESS",
	Interrupted:    "INTERRUPTED",
	Failed:         "FAILED",
}

// Formatted - returns status in human-readable format.
func (pi PitbullStatus) Formatted() string {
	return fmt.Sprintf("%s (%d)", pitbullStatusNames[pi], pi)
}

// ProgressInfo - helper struct describing Pitbull progress.
type ProgressInfo struct {
	Checked int64 `bson:"checked" json:"checked"`
	Total   int64 `bson:"total" json:"total"`
}

// Formatted - returns progress in human-readable format.
func (pi *ProgressInfo) Formatted() string {
	if pi == nil || pi.Total == 0 {
		return "0 / 0"
	}

	return fmt.Sprintf("%d / %d", pi.Checked, pi.Total)
}

// PitbullInstance - describes a single instance of PitbullInstance, backed by a host instance
// defined in ProviderInstance field.
type PitbullInstance struct {
	BaseModel `bson:",inline"`

	WalletString string `bson:"walletString" json:"walletString"`
	PasslistUrl  string `bson:"passlistUrl" json:"passlistUrl"`

	StartedAt   NullableTime `bson:"startedAt" json:"startedAt"`
	CompletedAt NullableTime `bson:"completedAt" json:"completedAt"`

	Status     PitbullStatus `bson:"status" json:"status"`
	Progress   *ProgressInfo `bson:"progress" json:"progress"`
	LastOutput string        `bson:"lastOutput" json:"lastOutput"`

	ProviderName string `bson:"providerName" json:"providerName"`

	HostInstance    host.HostInstance `bson:"-" json:"hostInstance"`
	HostInstanceRaw bson.Raw          `bson:"hostInstanceRaw" json:"-"`
}

type marshalablePitbullInstance PitbullInstance

// NewPitbullInstance - returns new PitbullInstance instance.
func NewPitbullInstance(host host.HostInstance, passlistUrl, walletString string) *PitbullInstance {
	instance := &PitbullInstance{
		Status:       WaitingForHost,
		PasslistUrl:  passlistUrl,
		WalletString: walletString,
		HostInstance: host,
	}

	instance.ID = primitive.NewObjectID()

	return instance
}

func (pi *PitbullInstance) SetHost(hostInstance host.HostInstance) {
	pi.HostInstance = hostInstance
	pi.ProviderName = hostInstance.ProviderName()
}

// Active - returns true if PitbullInstance is in one of the "active" states,
// false otherwise.
func (pi *PitbullInstance) Active() bool {
	return pi.Status == HostStarting ||
		pi.Status == Waiting ||
		pi.Status == Running
}

// Completed - returns true if instance's status is either Finished or Success.
func (pi *PitbullInstance) Completed() bool {
	return pi.Status == Finished || pi.Status == Success
}

// AllPasswordsChecked - returns true if all password have been checked, according to
// ProgressInfo.
func (pi *PitbullInstance) AllPasswordsChecked() bool {
	if pi.Progress == nil {
		return false
	}

	return pi.Progress.Total > 0 && (pi.Progress.Checked == pi.Progress.Total)
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
		// In this case, raw status will be "WAITING".
		pi.Status = Waiting
	}
}

// SetProgress - helper function for converting raw progress command output into
// ProgressInfo struct.
func (pi *PitbullInstance) SetProgress(rawProgress string) error {
	if pi.Progress == nil {
		pi.Progress = &ProgressInfo{
			Checked: 0,
			Total:   0,
		}
	}

	// If the command returned "NO_PROGRESS_AVAILABLE, and Pitbull is either WAITING or FINISHED,
	// we don't want to change the progress.
	if strings.Contains(rawProgress, "NO_PROGRESS_AVAILABLE") {
		return nil
	}

	parts := strings.Split(rawProgress, "of")

	if parts == nil || len(parts) < 2 || parts[0] == "" || parts[1] == "" {
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
