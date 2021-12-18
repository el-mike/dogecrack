package pitbull

// PitbullManager - entity responsible for managing Pitbull instances.
type PitbullManager interface {
	// Sync - checks and synchronizes provider's instances.
	Sync() ([]PitbullInstance, error)

	// RunInstance - runs single instance.
	RunInstance(fileUrl, wallet string) (PitbullInstance, error)

	// CheckInstance - returns the current status and output of Pitbull instance.
	CheckInstance(instance PitbullInstance) (PitbullStatus, string, error)
}
