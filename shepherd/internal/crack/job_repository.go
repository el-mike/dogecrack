package crack

import (
	"context"

	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const JobsCollection = "crack_jobs"

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
func (jr *JobRepository) GetById(id string) (*models.CrackJob, error) {
	collection := jr.db.Collection(JobsCollection)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectId}

	result := collection.FindOne(context.TODO(), filter)

	job := &models.CrackJob{}

	if err := result.Decode(job); err != nil {
		return nil, err
	}

	return job, nil
}

// RescheduleProcessingJobs - marks all "processing" jobs as "rescheduled".
func (jr *JobRepository) RescheduleProcessingJobs(jobIds []string) error {
	collection := jr.db.Collection(JobsCollection)

	objectIds := []primitive.ObjectID{}

	for _, jobId := range jobIds {
		objectId, err := primitive.ObjectIDFromHex(jobId)
		if err != nil {
			return err
		}

		objectIds = append(objectIds, objectId)
	}

	filter := bson.D{
		{"_id", bson.D{{"$in", objectIds}}},
	}

	update := bson.D{
		{"$set", bson.D{
			{"status", models.JobStatus.Rescheduled},
			{"lastScheduledAt", models.NullableTimeNow().Time},
		}},
	}

	if _, err := collection.UpdateMany(context.TODO(), filter, update); err != nil {
		return err
	}

	return nil
}

func (jr *JobRepository) GetAll(payload *models.CrackJobsListPayload) ([]*models.CrackJob, int, error) {
	collection := jr.db.Collection(JobsCollection)

	statuses := payload.Statuses
	pageSize := payload.PageSize
	page := payload.Page
	keyword := payload.Keyword
	passlistUrl := payload.PasslistUrl
	name := payload.Name

	// We ignore error return on purpose - when payload.JobId is incorrect,
	// jobId will be NilObjectID, and we can test against that.
	jobId, _ := primitive.ObjectIDFromHex(payload.JobId)

	statusesFilter := bson.D{}

	if len(statuses) > 0 {
		statusesFilter = bson.D{
			{"status", bson.D{{"$in", statuses}}},
		}
	}

	keywordFilter := bson.D{}

	if keyword != "" {
		keywordFilter = bson.D{
			{"keyword", primitive.Regex{Pattern: keyword, Options: ""}},
		}
	}

	passlistUrlFilter := bson.D{}

	if passlistUrl != "" {
		passlistUrlFilter = bson.D{
			{"passlistUrl", primitive.Regex{Pattern: passlistUrl, Options: ""}},
		}
	}

	jobIdFilter := bson.D{}

	if jobId != primitive.NilObjectID {
		jobIdFilter = bson.D{
			{"_id", jobId},
		}
	}

	nameFilter := bson.D{}

	if name != "" {
		nameFilter = bson.D{
			{"name", primitive.Regex{Pattern: name, Options: ""}},
		}
	}

	match := bson.D{{"$match", bson.D{
		{"$and", bson.A{statusesFilter, keywordFilter, passlistUrlFilter, jobIdFilter, nameFilter}},
	}}}

	lookup, unwind := jr.lookupAndUnwindInstance()

	sort := bson.D{{"$sort", bson.D{{"createdAt", -1}}}}

	// Page is 1-based, so we need to subtract one.
	skip := pageSize * (page - 1)
	limit := pageSize

	facet := bson.D{
		{"$facet", bson.D{
			{"pageInfo", bson.A{
				bson.D{{"$count", "total"}},
				bson.D{{"$addFields", bson.D{
					{"page", page},
					{"pageSize", pageSize},
				}}},
			}},
			{"data", bson.A{
				bson.D{{"$skip", skip}},
				bson.D{{"$limit", limit}},
				lookup,
				unwind,
			}},
		}},
	}

	unwindPageInfo := bson.D{
		{"$unwind", bson.D{{"path", "$pageInfo"}}},
	}

	pipeline := mongo.Pipeline{match, sort, facet, unwindPageInfo}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, 0, err
	}

	result := []*models.PagedCrackJobs{}

	if err = cursor.All(context.TODO(), &result); err != nil {
		return nil, 0, err
	}

	paged := models.NewPagedCrackJobs()
	// cursor.All() always returns a slice, therefore we need to get the first element.
	if len(result) > 0 {
		paged = result[0]
	}

	jobs := paged.Data

	if jobs == nil {
		jobs = []*models.CrackJob{}
	}

	return jobs, paged.PageInfo.Total, nil
}

// Create - saves a new PitbullJob to the DB.
func (jr *JobRepository) Create(job *models.CrackJob) error {
	collection := jr.db.Collection(JobsCollection)

	job.CreatedAt = models.NullableTimeNow()
	job.UpdatedAt = models.NullableTimeNow()

	result, err := collection.InsertOne(context.TODO(), job)
	if err != nil {
		return err
	}

	job.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

// Update - updates given PitbullJob in the DB.
func (jr *JobRepository) Update(job *models.CrackJob) error {
	collection := jr.db.Collection(JobsCollection)

	job.UpdatedAt = models.NullableTimeNow()

	payload := bson.D{{Key: "$set", Value: job}}

	if _, err := collection.UpdateByID(context.TODO(), job.ID, payload); err != nil {
		return err
	}

	return nil
}

// GetStatistics - returns statistics for CrackJobs.
func (jr *JobRepository) GetStatistics() (*models.CrackJobsStatistics, error) {
	collection := jr.db.Collection(JobsCollection)

	flattenCountDescription := func(field string) bson.E {
		// Resolves to { $arrayElemAt: ["$field.field", 0] }
		return bson.E{field, bson.D{{"$arrayElemAt", bson.A{"$" + field + "." + field, 0}}}}
	}

	facet := bson.D{
		{"$facet", bson.D{
			{"all", bson.A{
				bson.D{{"$count", "all"}},
			}},
			{"acknowledged", bson.A{
				bson.D{{"$match", bson.D{{"status", models.JobStatus.Acknowledged}}}},
				bson.D{{"$count", "acknowledged"}},
			}},
			{"processing", bson.A{
				bson.D{{"$match", bson.D{{"status", models.JobStatus.Processing}}}},
				bson.D{{"$count", "processing"}},
			}},
			{"queued", bson.A{
				bson.D{{"$match", bson.D{{"status", bson.D{{"$in", bson.A{models.JobStatus.Scheduled, models.JobStatus.Rescheduled}}}}}}},
				bson.D{{"$count", "queued"}},
			}},
			{"rejected", bson.A{
				bson.D{{"$match", bson.D{{"status", models.JobStatus.Rejected}}}},
				bson.D{{"$count", "rejected"}},
			}},
		}},
	}

	project := bson.D{
		{"$project", bson.D{
			flattenCountDescription("all"),
			flattenCountDescription("acknowledged"),
			flattenCountDescription("processing"),
			flattenCountDescription("queued"),
			flattenCountDescription("rejected"),
		}},
	}

	pipeline := mongo.Pipeline{facet, project}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	result := []*models.CrackJobsStatistics{}

	if err = cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	if result == nil || len(result) == 0 {
		return &models.CrackJobsStatistics{}, nil
	}

	return result[0], nil
}

func (jr *JobRepository) lookupAndUnwindInstance() (bson.D, bson.D) {
	lookup := bson.D{
		{"$lookup", bson.D{
			{"from", "pitbull_instances"},
			{"localField", "instanceId"},
			{"foreignField", "_id"},
			{"as", "instance"}},
		},
	}

	// Without $unwind stage, aggregation will return an array in "instance" field,
	// therefore making instance unmarshaling impossible and returning a Zero-value for the field.
	unwind := bson.D{
		{"$unwind", bson.D{{"path", "$instance"}, {"preserveNullAndEmptyArrays", true}}},
	}

	return lookup, unwind
}
