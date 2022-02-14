package auth

import (
	"context"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"github.com/go-redis/redis/v8"
)

var ctx = context.TODO()

// Cache - a cache for storing user session ids.
type Cache struct {
	redisClient *redis.Client

	expiration time.Duration
}

// NewCache - returns new Cache instance.
func NewCache(expiration time.Duration) *Cache {
	return &Cache{
		redisClient: persist.GetRedisClient(),

		expiration: expiration,
	}
}

// GetSessionId - gets a userId by given sessionId if exists.
func (cc *Cache) GetUserBySessionId(sessionId string) (string, error) {
	result, err := cc.redisClient.Get(ctx, sessionId).Result()
	if err == redis.Nil {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return result, nil
}

// SaveSessionId - saves a session id for given user.
func (cc *Cache) SaveSessionId(sessionId, userId string) error {
	_, err := cc.redisClient.SetEX(ctx, sessionId, userId, cc.expiration).Result()
	if err != nil {
		return err
	}

	return nil
}

// DeleteSessionId - deletes given sessionId.
func (cc *Cache) DeleteSessionId(sessionId string) error {
	_, err := cc.redisClient.Del(ctx, sessionId).Result()
	if err != nil {
		return err
	}

	return nil
}
