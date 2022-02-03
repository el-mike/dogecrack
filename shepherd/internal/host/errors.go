package host

import "fmt"

// HostInstanceNotFound - error thrown when instance with given ID was not found
// by provider.
type HostInstanceNotFound struct {
	InstanceId int
}

// NewHostInstanceNotFound - returns new InstanceNotFound instance.
func NewHostInstanceNotFound(instanceId int) *HostInstanceNotFound {
	return &HostInstanceNotFound{
		InstanceId: instanceId,
	}
}

// Error - Error implementation.
func (e *HostInstanceNotFound) Error() string {
	return fmt.Sprintf("Instance with ID: %d not found", e.InstanceId)
}
