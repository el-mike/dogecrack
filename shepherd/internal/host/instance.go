package host

import "fmt"

// HostStatus - enum describing the status of the Pitbull's host machine.
type HostStatus int8

const (
	Starting HostStatus = iota
	Running
	Destroyed
	Unknown
)

var hostStatusNames = map[HostStatus]string{
	Starting:  "STARTING",
	Running:   "RUNNING",
	Destroyed: "DESTROYED",
	Unknown:   "UNKNOWN",
}

// Formatted - returns HostStatus in human-readable format.
func (hs HostStatus) Formatted() string {
	return fmt.Sprintf("%s (%d)", hostStatusNames[hs], hs)
}

// HostInstance - representation of a single machine instance as returned
// by given provider.
type HostInstance interface {
	// ProviderName - returns provider's name.
	ProviderName() string
	// ProviderId - id of actual machine/instance of given provider.
	ProviderId() int
	// HostStatus - returns status value as defined internally.
	HostStatus() HostStatus
}
