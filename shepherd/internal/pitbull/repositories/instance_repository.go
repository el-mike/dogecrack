package repositories

import (
	"context"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const instancesCollection = "instances"

// InstanceRepository - MongoDB-backed repository for handling Pitbull instances.
type InstanceRepository struct {
	db *mongo.Database
}

// NewInstanceRepository - returns new InstanceRepository.
func NewInstanceRepository() *InstanceRepository {
	return &InstanceRepository{
		db: persist.GetDatabase(),
	}
}

// GetInstanceById - returns an instance with given id.
func (ir *InstanceRepository) GetInstanceById(id string) (*models.PitbullInstance, error) {
	collection := ir.db.Collection(instancesCollection)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectId}

	result := collection.FindOne(context.TODO(), filter)

	instance := &models.PitbullInstance{}

	if err := result.Decode(instance); err != nil {
		return nil, err
	}

	return instance, nil
}

// CreateInstance - saves a new Pitbull instance to the DB.
func (ir *InstanceRepository) CreateInstance(pitbull *models.PitbullInstance) error {
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

func (ir *InstanceRepository) UpdateInstance(pitbull *models.PitbullInstance) error {
	collection := ir.db.Collection(instancesCollection)

	pitbull.UpdatedAt = time.Now()

	payload := bson.D{{Key: "$set", Value: pitbull}}

	if _, err := collection.UpdateByID(context.TODO(), pitbull.ID, payload); err != nil {
		return err
	}

	return nil
}

// GetActiveInstances - returns all instances that are active, i.e. thepir status is different
// than FINISHED (4).
func (ir *InstanceRepository) GetActiveInstances() ([]*models.PitbullInstance, error) {
	collection := ir.db.Collection(instancesCollection)

	filter := bson.D{
		{"status", bson.D{{"$nin", bson.A{models.Finished, models.Success}}}},
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