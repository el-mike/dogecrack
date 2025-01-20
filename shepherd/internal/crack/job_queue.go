package crack

import (
	"context"

	"github.com/el-mike/dogecrack/shepherd/internal/persist"
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
func NewJobQueue() *JobQueue {
	return &JobQueue{
		redisClient:     persist.GetRedisClient(),
		waitingQueue:    waitingQueue,
		processingQueue: processingQueue,
	}
}

// QueueRun - adds specified jobIds to the waiting queue.
func (jq *JobQueue) Enqueue(jobIds ...string) error {
	_, err := jq.redisClient.LPush(ctx, jq.waitingQueue, jobIds).Result()

	return err
}

// Dequeue - pops a single jobId from the waiting queue, and pushes it to processing queue.
// This operation is atomic.
func (jq *JobQueue) Dequeue() (string, error) {
	// We are using RPOPLPUSH, so we can retry jobs that failed before finishing.
	jobId, err := jq.redisClient.RPopLPush(ctx, jq.waitingQueue, jq.processingQueue).Result()
	// If error is "redis.Nil", that means list was empty, but we don't want to
	// treat it as error - therefore, we use nil, nil.
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return jobId, nil
}

// Ack - acknowledges a single job and removes it from the processing queue.
func (jq *JobQueue) Ack(jobId string) error {
	return jq.removeProcessing(jobId)
}

// Reschedule - moves a single job from the processing queue to working queue.
func (jq *JobQueue) Reschedule(jobId string) error {
	if err := jq.removeProcessing(jobId); err != nil {
		return err
	}

	return jq.Enqueue(jobId)
}

// Reject - rejects a single job and removes it from any queue.
func (jq *JobQueue) Reject(jobId string) error {
	err := jq.removeProcessing(jobId)
	if err != nil {
		return err
	}

	_, err = jq.redisClient.LRem(ctx, jq.waitingQueue, 1, jobId).Result()

	return err
}

// ClearProcessing - removes all
func (jq *JobQueue) RescheduleAllProcessing() ([]string, error) {
	jobsIds, err := jq.GetProcessing()

	if err != nil {
		return nil, err
	}

	if len(jobsIds) == 0 {
		return nil, nil
	}

	if err := jq.Enqueue(jobsIds...); err != nil {
		return nil, err
	}

	_, err = jq.redisClient.Del(ctx, jq.processingQueue).Result()

	return jobsIds, err
}

// GetProcessing - returns all IDs in processingQueue.
func (jq *JobQueue) GetProcessing() ([]string, error) {
	jobsIds, err := jq.redisClient.LRange(ctx, jq.processingQueue, 0, -1).Result()
	// If the err is redis.Nil, it just means the queue is empty.
	if err == redis.Nil {
		return nil, nil
	}

	return jobsIds, err
}

// removeProcessing - removes a single job id from processing queue.
func (jq *JobQueue) removeProcessing(jobId string) error {
	_, err := jq.redisClient.LRem(ctx, jq.processingQueue, 1, jobId).Result()

	return err
}
