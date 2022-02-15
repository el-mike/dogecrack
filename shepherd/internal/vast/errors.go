package vast

import "fmt"

// VastCLIError - error returned from vast CLI command.
type VastCLIError struct {
	StdoutLog string
	StderrLog string
}

// NewVastCLIError - returns new VastCLIError instance.
func NewVastCLIError(stdoutLog, stderrLog string) *VastCLIError {
	return &VastCLIError{
		StdoutLog: stdoutLog,
		StderrLog: stderrLog,
	}
}

// Error - error implementation.
func (e *VastCLIError) Error() string {
	return fmt.Sprintf("Vast CLI error occured.\n[Stdout]: %s\n[Stderr]: %s\n", e.StdoutLog, e.StderrLog)
}
