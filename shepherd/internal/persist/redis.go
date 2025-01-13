package persist

import (
	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// GetRedisClient - singleton implementation for redisClient.
func GetRedisClient() *redis.Client {
	return redisClient
}

// InitRedis - inits Redis instance connection.
func InitRedis(connectionString string) *redis.Client {
	if redisClient != nil {
		return redisClient
	}

	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		panic(err)
	}
	
	redisClient = redis.NewClient(opt)

	return redisClient
}
