package common

type InstanceStatus int8

const (
	Starting InstanceStatus = iota
	Running
	Interrupted
	Failed
	Finished
	Success
)

type PitbullInstance struct {
	ID               string
	Name             string
	Rules            []string
	ProviderInstance ProviderInstance
	Status           InstanceStatus
}
