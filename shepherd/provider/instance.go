package provider

// InstanceStatus - common status enum.
type InstanceStatus int8

const (
	Starting InstanceStatus = iota
	Running
	Interrupted
	Failed
	Finished
	Success
)

// ProviderInstance - representation of a single machine instance as returned
// by given provider.
type ProviderInstance interface {
	// ProviderName - returns provider's name.
	ProviderName() string
	// ProviderId - id of actual machine/instance of given provider.
	ProviderId() int
	// PitbullStatus - returns status value as defined internally.
	PitbullStatus() InstanceStatus
}
