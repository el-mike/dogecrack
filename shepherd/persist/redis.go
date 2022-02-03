package persist

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// GetRedisClient - singleton implementation for redisClient.
func GetRedisClient() *redis.Client {
	return redisClient
}

// InitRedis - inits Redis instance connection.
func InitRedis(ctx context.Context, host, port string) *redis.Client {
	if redisClient != nil {
		return redisClient
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
	})

	return redisClient
}
