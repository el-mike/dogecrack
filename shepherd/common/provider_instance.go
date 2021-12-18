package common

// ProviderInstance - representation of a single machine instance as returned
// by given provider.
type ProviderInstance interface {
	// ProviderId - id of actual machine/instance of given provider.
	ProviderId() int
	// PitbullId - id for internal use, set by Shepherd application.
	PitbullId() string
}
