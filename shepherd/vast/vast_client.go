package vast

type VastClient interface {
	// GetInstances - returns current instances.
	GetInstances() ([]*VastInstance, error)

	// StartInstance - starts new Vast.ai instance. Waits for starting process to be over.
	StartInstance(offerId int) (*VastCreateResponse, error)

	// DestroyInstance - stops a Vast.ai instance with given id.
	DestroyInstance(instandeId int) error

	// GetInstance - returns single, existing (rented) instance based on passed id.
	GetInstance(instanceId int) (*VastInstance, error)

	// GetOfferByCriteria - returns first offer matching the criteria.
	GetOfferByCriteria(criteria string) (*VastOffer, error)
}
