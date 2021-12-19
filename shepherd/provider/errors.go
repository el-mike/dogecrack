package provider

import "fmt"

// InstanceNotFound - error thrown when instance with given ID was not found
// by provider.
type InstanceNotFound struct {
	InstanceId int
}

// NewInstanceNotFound - returns new InstanceNotFound instance.
func NewInstanceNotFound(instanceId int) *InstanceNotFound {
	return &InstanceNotFound{
		InstanceId: instanceId,
	}
}

// Error - Error implementation.
func (e *InstanceNotFound) Error() string {
	return fmt.Sprintf("Instance with ID: %d not found", e.InstanceId)
}
