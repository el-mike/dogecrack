package provider

// ProviderInstance - representation of a single machine instance as returned
// by given provider.
type ProviderInstance interface {
	// ProviderName - returns provider's name.
	ProviderName() string
	// ProviderId - id of actual machine/instance of given provider.
	ProviderId() int
	// PitbullId - id for internal use, set by Shepherd application.
	PitbullId() string
	// PitbullStatus - returns status value as defined internally.
	PitbullStatus() InstanceStatus
}
