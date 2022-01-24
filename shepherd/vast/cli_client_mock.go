package vast

type VastCLIClientMock struct{}

func NewVastCLIClientMock() *VastCLIClientMock {
	return &VastCLIClientMock{}
}

// GetInstances - returns current instances.
func (vc *VastCLIClientMock) GetInstances() ([]*VastInstance, error) {
	return nil, nil
}

// StartInstance - starts new Vast.ai instance. Waits for starting process to be over.
func (vc *VastCLIClientMock) StartInstance(offerId int) (*VastCreateResponse, error) {
	return nil, nil
}

// GetInstance - returns single, existing (rented) instance based on passed id.
func (vc *VastCLIClientMock) GetInstance(instanceId int) (*VastInstance, error) {
	return nil, nil
}

// GetOfferByCriteria - returns first offer matching the criteria.
func (vc *VastCLIClientMock) GetOfferByCriteria(criteria string) (*VastOffer, error) {
	return nil, nil
}
