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
	SkipCount  int64             `bson:"skipCount" json:"skipCount"`
}

func NewPitbull(skipCount int64) *Pitbull {
	return &Pitbull{
		Status:    PitbullStatus.Waiting,
		SkipCount: skipCount,
	}
}

// AllPasswordsChecked - returns true if all password have been checked, based on ProgressInfo.
func (pl *Pitbull) AllPasswordsChecked() bool {
	if pl.Progress == nil {
		return false
	}

	return pl.Progress.Total > 0 && (pl.Progress.Checked == pl.Progress.Total)
}

// AnyPasswordsChecked - returns true if at least one password has been checked, based on ProgressInfo.
func (pl *Pitbull) AnyPasswordsChecked() bool {
	if pl.Progress == nil {
		return false
	}

	return pl.Progress.Checked > 0
}

// AnyPasswordsSkipped - returns true if at least one password has been skipped.
func (pl *Pitbull) AnyPasswordsSkipped() bool {
	if pl.Progress == nil {
		return false
	}

	return pl.SkipCount > 0
}

// GetResumeCount - returns the number of password Pitbull could skip if rerun.
func (pl *Pitbull) GetResumeCount() int64 {
	if pl.Progress == nil {
		return 0
	}

	// If given Pitbull was already run with SkipCount, we need to accommodate it as well for the resume count.
	return pl.SkipCount + pl.Progress.Checked
}

// Done - returns true if status is either Finished or Success.
func (pl *Pitbull) Done() bool {
	return pl.Status == PitbullStatus.Finished || pl.Status == PitbullStatus.Success
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

// GetOutputLastLine - returns last line of current available Pitbull output.
func (pl *Pitbull) GetOutputLastLine() string {
	lines := strings.Split(pl.LastOutput, "\n")

	if len(lines) == 0 {
		return "<NO_OUTPUT_AVAILABLE>"
	}

	// Last line can simply be an empty line - therefore, we iterate from the end
	// to get last non-empty line.
	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i] != "" {
			return lines[i]
		}
	}

	return ""
}
