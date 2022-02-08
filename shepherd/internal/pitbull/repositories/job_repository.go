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

const jobsCollection = "jobs"

// JobRepository - MongoDB-backed repository for handling Pitbull jobs.
type JobRepository struct {
	db *mongo.Database
}

// NewJobRepository - returns new PitbullJobRepository instance.
func NewJobRepository() *JobRepository {
	return &JobRepository{
		db: persist.GetDatabase(),
	}
}

// GetById - returns a single Job with given id.
func (jr *JobRepository) GetById(id string) (*models.PitbullJob, error) {
	collection := jr.db.Collection(jobsCollection)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectId}

	result := collection.FindOne(context.TODO(), filter)

	job := &models.PitbullJob{}

	if err := result.Decode(job); err != nil {
		return nil, err
	}

	return job, nil
}

// Create - saves a new PitbullJob to the DB.
func (jr *JobRepository) Create(job *models.PitbullJob) error {
	collection := jr.db.Collection(jobsCollection)

	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()

	result, err := collection.InsertOne(context.TODO(), job)
	if err != nil {
		return err
	}

	job.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

// Update - updates given PitbullJob in the DB.
func (jr *JobRepository) Update(job *models.PitbullJob) error {
	collection := jr.db.Collection(jobsCollection)

	job.UpdatedAt = time.Now()

	payload := bson.D{{Key: "$set", Value: job}}

	if _, err := collection.UpdateByID(context.TODO(), job.ID, payload); err != nil {
		return err
	}

	return nil
}

func (jr *JobRepository) GetAll(statuses []models.JobStatus) ([]*models.PitbullJob, error) {
	collection := jr.db.Collection(jobsCollection)

	filter := bson.D{}

	if len(statuses) > 0 {
		filter = bson.D{
			{"status", bson.D{{"$in", statuses}}},
		}
	}

	match := bson.D{{"$match", filter}}
	lookup := bson.D{
		{"$lookup", bson.D{
			{"from", "instances"},
			{"localField", "instanceId"},
			{"foreignField", "_id"},
			{"as", "instance"}},
		},
	}

	// Without $unwind stage, aggregation will return an array in "instance" field,
	// therefore making instance unmarshaling impossible and returning a Zero-value for the field.
	unwind := bson.D{
		{"$unwind", bson.D{{"path", "$instance"}}},
	}

	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{match, lookup, unwind})
	if err != nil {
		return nil, err
	}

	var jobs []*models.PitbullJob

	if err = cursor.All(context.TODO(), &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}
