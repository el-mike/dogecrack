package persist

import (
	"github.com/el-mike/dogecrack/shepherd/models"
	"go.mongodb.org/mongo-driver/mongo"
)

const instancesCollection = "instances"

type InstanceRepository struct {
	db *mongo.Database
}

func NewInstanceRepository() *InstanceRepository {
	return &InstanceRepository{
		db: GetDatabase(),
	}
}

func (ir *InstanceRepository) GetByProviderId() {}

func (ir *InstanceRepository) GetActiveInstances() ([]*models.PitbullInstance, error) {
	// collection := ir.db.Collection(instancesCollection)

	return nil, nil
}
