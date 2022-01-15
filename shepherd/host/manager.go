package host

// HostManager - entity responsible for managing actual, provider's instances.
type HostManager interface {
	// Sync - checks and synchronizes host's instances.
	Sync() ([]HostInstance, error)

	// RunInstance - runs single instance.
	RunInstance(fileUrl, wallet string) (HostInstance, error)

	// GetInstance - returns active(rented) instance with given ID.
	GetInstance(instanceId int) (HostInstance, error)
}
