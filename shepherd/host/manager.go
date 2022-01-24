package host

// HostManager - entity responsible for managing actual, provider's instances and
// underlying Pitbull process.
type HostManager interface {
	// Sync - checks and synchronizes host's instances.
	Sync() ([]HostInstance, error)

	// RunInstance - runs single instance.
	RunInstance(fileUrl, wallet string) (HostInstance, error)

	// GetInstance - returns active(rented) instance with given ID.
	GetInstance(instanceId int) (HostInstance, error)

	// GetPitbullStatus - returns Pitbull's status command output for a host instance
	// with given ID.
	GetPitbullStatus(instance HostInstance) (string, error)

	// GetPitbullProgress - returns Pitbull's progress command output for a host instance
	// with given ID.
	GetPitbullProgress(instance HostInstance) (string, error)

	//RunDirectCommand - runs given command directly on a host instance.
	RunDirectCommand(instance HostInstance, cmd string) (string, error)
}
