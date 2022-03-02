package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// PitbullStatusEnum - enum describing the status of a Pitbull process.
type PitbullStatusEnum int8

// PitbullStatus - "namespaced" struct containing PitbullStatus values.
var PitbullStatus = struct {
	Waiting  PitbullStatusEnum
	Running  PitbullStatusEnum
	Finished PitbullStatusEnum
	Success  PitbullStatusEnum
}{
	Waiting:  0,
	Running:  1,
	Finished: 2,
	Success:  3,
}

var pitbullStatusByName = map[PitbullStatusEnum]string{
	PitbullStatus.Waiting:  "WAITING",
	PitbullStatus.Running:  "RUNNING",
	PitbullStatus.Finished: "FINISHED",
	PitbullStatus.Success:  "SUCCESS",
}

// Formatted - returns status in human-readable format.
func (pi PitbullStatusEnum) Formatted() string {
	return fmt.Sprintf("%s (%d)", pitbullStatusByName[pi], pi)
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

// Pitbull - describes a single Pitbull process.
type Pitbull struct {
	Status     PitbullStatusEnum `bson:"status" json:"status"`
	Progress   *ProgressInfo     `bson:"progress" json:"progress"`
	LastOutput string            `bson:"lastOutput" json:"lastOutput"`
}

func NewPitbull() *Pitbull {
	return &Pitbull{
		Status: PitbullStatus.Waiting,
	}
}

// AllPasswordsChecked - returns true if all password have been checked, according to
// ProgressInfo.
func (pl *Pitbull) AllPasswordsChecked() bool {
	if pl.Progress == nil {
		return false
	}

	return pl.Progress.Total > 0 && (pl.Progress.Checked == pl.Progress.Total)
}

// ParseRawStatus - helper function for converting raw status command output into
// PitbullStatus enum.
func (pl *Pitbull) ParseRawStatus(rawStatus string) {
	if strings.Contains(rawStatus, "SUCCESS") {
		pl.Status = PitbullStatus.Success
	} else if strings.Contains(rawStatus, "RUNNING") {
		pl.Status = PitbullStatus.Running
	} else if strings.Contains(rawStatus, "FINISHED") {
		pl.Status = PitbullStatus.Finished
	} else {
		// In this case, raw status will be "WAITING".
		pl.Status = PitbullStatus.Waiting
	}
}

// ParseProgress - helper function for converting raw progress command output into
// ProgressInfo struct.
func (pl *Pitbull) ParseProgress(rawProgress string) error {
	if pl.Progress == nil {
		pl.Progress = &ProgressInfo{
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

	pl.Progress.Checked = int64(checked)
	pl.Progress.Total = int64(total)

	return nil
}
