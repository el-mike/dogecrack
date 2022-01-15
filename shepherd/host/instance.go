package host

// HostStatus - enum describing the status of the Pitbull's host machine.
type HostStatus int8

const (
	Starting HostStatus = iota
	Running
	Finished
)

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
