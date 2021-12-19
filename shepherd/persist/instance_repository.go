package persist

import (
	"context"
	"time"

	"github.com/el-mike/dogecrack/shepherd/models"
	"github.com/el-mike/dogecrack/shepherd/provider"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const instancesCollection = "instances"

// InstanceRepository - MongoDB repository for handling Pitbull instances.
type InstanceRepository struct {
	db *mongo.Database
}

// NewInstanceRepository - returns new InstanceRepository.
func NewInstanceRepository() *InstanceRepository {
	return &InstanceRepository{
		db: GetDatabase(),
	}
}

// SaveInstance - saves a new Pitbull instance to the DB.
func (ir *InstanceRepository) SaveInstance(pitbull *models.PitbullInstance) error {
	collection := ir.db.Collection(instancesCollection)

	pitbull.CreatedAt = time.Now()
	pitbull.UpdatedAt = time.Now()

	result, err := collection.InsertOne(context.TODO(), pitbull)
	if err != nil {
		return err
	}

	pitbull.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

func (ir *InstanceRepository) GetByProviderId() {}

func (ir *InstanceRepository) GetActiveInstances() ([]*models.PitbullInstance, error) {
	collection := ir.db.Collection(instancesCollection)

	filter := bson.D{
		{"status", bson.D{{"$ne", provider.Finished}}},
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var results []*models.PitbullInstance

	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
