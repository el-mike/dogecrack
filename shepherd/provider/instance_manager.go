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

// ProviderInstanceManager - entity responsible for managing actual, provider's instances.
type ProviderInstanceManager interface {
	// Sync - checks and synchronizes provider's instances.
	Sync() ([]ProviderInstance, error)

	// RunInstance - runs single instance.
	RunInstance(fileUrl, wallet string) (ProviderInstance, error)

	// CheckInstance - returns the current status and output of Pitbull instance.
	CheckInstance(instance ProviderInstance) (InstanceStatus, string, error)
}
