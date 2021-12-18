package pitbull

type PitbullStatus int8

const (
	Starting PitbullStatus = iota
	Running
	Interrupted
	Failed
	Finished
)

// PitbullInstance - representation of a single Pitbull instance.
type PitbullInstance interface {
	// PitbullId - id for internal use, set by Shepherd application.
	PitbullId() string
	// ProviderId - id of actual machine/instance of given provider.
	ProviderId() string
	// PitbullStatus - current instance's status as defined in Shepherd application.
	PitbullStatus() PitbullStatus
}

type Pitbull struct {
	ID         string
	Name       string
	Rules      []string
	ProviderId string
	Status     PitbullStatus
}
