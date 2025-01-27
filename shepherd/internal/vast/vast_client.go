package vast

import "github.com/el-mike/dogecrack/shepherd/internal/vast/models"

type VastClient interface {
	// GetInstances - returns current instances.
	GetInstances() ([]*models.Instance, error)

	// StartInstance - starts new Vast.ai instance. Waits for starting process to be over.
	StartInstance(offerId int) (*models.CreateResponse, error)

	// DestroyInstance - stops a Vast.ai instance with given id.
	DestroyInstance(instanceId int) error

	// GetInstance - returns single, existing (rented) instance based on passed id.
	GetInstance(instanceId int) (*models.Instance, error)

	// GetOfferByCriteria - returns first offer matching the criteria.
	GetOfferByCriteria(criteria string) (*models.Offer, error)

	// GetAvailableCredit - returns account's remaining credits.
	GetAvailableCredit() (float64, error)
}
