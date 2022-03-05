package host

import "fmt"

// NoHostMachinesAvailable - thrown when there are no machines that could be retned.
type NoHostMachinesAvailable struct{}

// NoHostMachinesAvailable - returns new NoHostMachinesAvailable instance.
func NewNoHostMachinesAvailable() *NoHostMachinesAvailable {
	return &NoHostMachinesAvailable{}
}

// Error - Error implementation.
func (e *NoHostMachinesAvailable) Error() string {
	return fmt.Sprintf("No host machines are available to be rented")
}

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

// HostInstanceNotAvailable - error thrown when instance with given ID could not
// be reached for some reason.
type HostInstanceNotAvailable struct {
	InstanceId int
}

// NewHostInstanceNotAvailable - returns new HostInstanceNotAvailable instance.
func NewHostInstanceNotAvailable(instanceId int) *HostInstanceNotAvailable {
	return &HostInstanceNotAvailable{
		InstanceId: instanceId,
	}
}

func (e *HostInstanceNotAvailable) Error() string {
	return fmt.Sprintf("Instance with ID %d is not available right now", e.InstanceId)
}
