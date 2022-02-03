package persist

import (
	"context"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const instancesCollection = "instances"

// PitbullInstanceRepository - MongoDB-backed repository for handling Pitbull instances.
type PitbullInstanceRepository struct {
	db *mongo.Database
}

// NewPitbullInstanceRepository - returns new InstanceRepository.
func NewPitbullInstanceRepository() *PitbullInstanceRepository {
	return &PitbullInstanceRepository{
		db: GetDatabase(),
	}
}

// GetInstanceById - returns an instance with given id.
func (pir *PitbullInstanceRepository) GetInstanceById(id string) (*models.PitbullInstance, error) {
	collection := pir.db.Collection(instancesCollection)

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
func (pir *PitbullInstanceRepository) CreateInstance(pitbull *models.PitbullInstance) error {
	collection := pir.db.Collection(instancesCollection)

	pitbull.CreatedAt = time.Now()
	pitbull.UpdatedAt = time.Now()

	result, err := collection.InsertOne(context.TODO(), pitbull)
	if err != nil {
		return err
	}

	pitbull.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

func (pir *PitbullInstanceRepository) UpdateInstance(pitbull *models.PitbullInstance) error {
	collection := pir.db.Collection(instancesCollection)

	pitbull.UpdatedAt = time.Now()

	updatePayload := bson.D{{Key: "$set", Value: pitbull}}

	_, err := collection.UpdateByID(context.TODO(), pitbull.ID, updatePayload)
	if err != nil {
		return err
	}

	return nil
}

// GetActiveInstances - returns all instances that are active, i.e. thepir status is different
// than FINISHED (4).
func (pir *PitbullInstanceRepository) GetActiveInstances() ([]*models.PitbullInstance, error) {
	collection := pir.db.Collection(instancesCollection)

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