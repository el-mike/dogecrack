package host

// HostManager - entity responsible for managing actual, provider's instances and
// underlying Pitbull process.
type HostManager interface {
	// CreateInstance - returns new host instance.
	CreateInstance() HostInstance

	// CanRunInstance - returns nil if new HostInstance can be run, error otherwise.
	CanRunInstance() error

	// RunInstance - runs single instance.
	RunInstance() (HostInstance, error)

	// GetInstance - returns active (rented) instance with given id.
	GetInstance(instanceId int) (HostInstance, error)

	// DestroyInstance - stops an instance with given ID.
	DestroyInstance(instanceId int) error

	// RunPitbullForPasslist - runs Pitbull process on given host instance with password provided under given url.
	RunPitbullForPasslist(
		instance HostInstance,
		walletString,
		passlistUrl string,
		skipCount,
		minLength,
		maxLength int64,
	) error

	// RunPitbullForTokenlist - runs Pitbull process on given host instance with provided tokenlist.
	RunPitbullForTokenlist(
		instance HostInstance,
		walletString,
		tokenlist string,
		skipCount,
		minLength,
		maxLength int64,
	) error

	// GetPitbullStatus - returns Pitbull's status command output for a host instance
	// with given id.
	GetPitbullStatus(instance HostInstance) (string, error)

	// GetPitbullProgress - returns Pitbull's progress command output for a host instance
	// with given id.
	GetPitbullProgress(instance HostInstance) (string, error)

	// GetPitbullOutput - returns Pitbull's output as returned by "pitbull output" command
	// for a host instance with given id.
	GetPitbullOutput(instance HostInstance) (string, error)

	//RunDirectCommand - runs given command directly on a host instance.
	RunDirectCommand(instance HostInstance, cmd string) (string, error)
}
