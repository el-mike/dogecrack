package pitbull

import (
	"context"
	"encoding/json"

	"github.com/el-mike/dogecrack/shepherd/internal/models"
	"github.com/go-redis/redis/v8"
)

const waitingQueue = "waitingQueue"
const processingQueue = "processingQueue"

var ctx = context.TODO()

// JobQueue - redis-based queue for scheduling Pitbull runs.
type JobQueue struct {
	redisClient *redis.Client

	waitingQueue    string
	processingQueue string
}

// NewJobQueue - returns new RunQueue instance.
func NewJobQueue(client *redis.Client) *JobQueue {
	return &JobQueue{
		redisClient:     client,
		waitingQueue:    waitingQueue,
		processingQueue: processingQueue,
	}
}

// QueueRun - adds a single PitbullRun to the job queue.
func (jq *JobQueue) Enqueue(job *models.PitbullJob) error {
	jobRaw, err := json.Marshal(job)
	if err != nil {
		return err
	}

	_, err = jq.redisClient.LPush(ctx, jq.waitingQueue, jobRaw).Result()

	return err
}

// Dequeue - pops a single PitbullRun from the job queue.
func (jq *JobQueue) Dequeue() (*models.PitbullJob, error) {
	// We are using RPOPLPUSH, so we can retry jobs that failed before finishing.
	jobRaw, err := jq.redisClient.RPopLPush(ctx, jq.waitingQueue, jq.processingQueue).Result()
	// If error is "redis.Nil", that means list was empty, but we don't want to
	// treat it as error - therefore, we use nil, nil.
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var job *models.PitbullJob
	if err := json.Unmarshal([]byte(jobRaw), &job); err != nil {
		return nil, err
	}

	return job, nil
}

func (jq *JobQueue) Ack(instanceId string) {

}
