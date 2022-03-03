package pitbull

import (
	"context"

	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const InstancesCollection = "pitbull_instances"

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
	collection := ir.db.Collection(InstancesCollection)

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
	collection := ir.db.Collection(InstancesCollection)

	pitbull.CreatedAt = models.NullableTimeNow()
	pitbull.UpdatedAt = models.NullableTimeNow()

	result, err := collection.InsertOne(context.TODO(), pitbull)
	if err != nil {
		return err
	}

	pitbull.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

func (ir *InstanceRepository) UpdateInstance(pitbull *models.PitbullInstance) error {
	collection := ir.db.Collection(InstancesCollection)

	pitbull.UpdatedAt = models.NullableTimeNow()

	payload := bson.D{{Key: "$set", Value: pitbull}}

	if _, err := collection.UpdateByID(context.TODO(), pitbull.ID, payload); err != nil {
		return err
	}

	return nil
}

// GetActiveInstances - returns all instances that are active, i.e. thepir status is different
// than FINISHED (4).
func (ir *InstanceRepository) GetActiveInstances() ([]*models.PitbullInstance, error) {
	collection := ir.db.Collection(InstancesCollection)

	filter := bson.D{
		{"status", bson.D{{"$nin", bson.A{models.PitbullStatus.Finished, models.PitbullStatus.Success}}}},
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

// GetOrphanInstances - returns "orphan" instances. Orphan instance is an instance that
// has no job assigned (meaning the job has been assiged with different instance), or that
// has one of the "active" statuses, but its job is already rejected/acknowledged.
func (ir *InstanceRepository) GetOrphanInstances() ([]*models.PitbullInstance, error) {
	collection := ir.db.Collection(InstancesCollection)

	lookup := bson.D{
		{"$lookup", bson.D{
			{"from", "crack_jobs"},
			{"localField", "_id"},
			{"foreignField", "instanceId"},
			{"as", "job"},
		}},
	}

	unwind := bson.D{
		{"$unwind", bson.D{
			{"path", "$job"},
			{"preserveNullAndEmptyArrays", true},
		}},
	}

	match := bson.D{
		{"$match", bson.D{
			{"status", bson.D{
				{"$in", bson.A{models.PitbullInstanceStatus.HostStarting, models.PitbullInstanceStatus.Running}},
			}},
			{"$or", bson.A{
				bson.D{{"job", nil}},
				bson.D{{"job.status", bson.D{
					{"$in", bson.A{models.JobStatus.Rejected, models.JobStatus.Acknowledged}},
				}}},
			}},
		}},
	}

	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{lookup, unwind, match})
	if err != nil {
		return nil, err
	}

	var instances []*models.PitbullInstance

	if err = cursor.All(context.TODO(), &instances); err != nil {
		return nil, err
	}

	return instances, nil
}

// GetStatistics - returns PitbullInstance statistics.
func (ir *InstanceRepository) GetStatistics() (*models.PitbullInstancesStatistics, error) {
	collection := ir.db.Collection(InstancesCollection)

	byStatusDescription := func(field string, status models.PitbullInstanceStatusEnum) bson.E {
		return bson.E{
			field, bson.A{
				bson.D{{"$match", bson.D{{"status", status}}}},
				bson.D{{"$count", field}},
			},
		}
	}

	flattenCountDescription := func(field string) bson.E {
		// Resolves to { $arrayElemAt: ["$field.field", 0] }
		return bson.E{field, bson.D{{"$arrayElemAt", bson.A{"$" + field + "." + field, 0}}}}
	}

	facet := bson.D{
		{"$facet", bson.D{
			{"all", bson.A{
				bson.D{{"$count", "all"}},
			}},
			byStatusDescription("waitingForHost", models.PitbullInstanceStatus.WaitingForHost),
			byStatusDescription("hostStarting", models.PitbullInstanceStatus.HostStarting),
			byStatusDescription("running", models.PitbullInstanceStatus.Running),
			byStatusDescription("completed", models.PitbullInstanceStatus.Completed),
			byStatusDescription("failed", models.PitbullInstanceStatus.Failed),
			byStatusDescription("interrupted", models.PitbullInstanceStatus.Interrupted),
			byStatusDescription("success", models.PitbullInstanceStatus.Success),
		}},
	}

	project := bson.D{
		{"$project", bson.D{
			flattenCountDescription("all"),
			flattenCountDescription("waitingForHost"),
			flattenCountDescription("hostStarting"),
			flattenCountDescription("running"),
			flattenCountDescription("completed"),
			flattenCountDescription("failed"),
			flattenCountDescription("interrupted"),
			flattenCountDescription("success"),
		}},
	}

	pipeline := mongo.Pipeline{facet, project}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	result := []*models.PitbullInstancesStatistics{}

	if err = cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	if result == nil || len(result) == 0 {
		return &models.PitbullInstancesStatistics{}, nil
	}

	return result[0], nil
}
