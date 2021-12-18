package pitbull

// PitbullManager - entity responsible for managing Pitbull instances.
type PitbullManager interface {
	// RunInstance - runs single instance.
	RunInstance(fileUrl, wallet string) (PitbullInstance, error)
}
