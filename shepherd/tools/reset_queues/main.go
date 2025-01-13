package main

import (
	"context"
	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"path/filepath"
)

func main() {
	rootPath, err := filepath.Abs("../../")
	if err != nil {
		panic(err)
	}

	appConfig, err := config.NewAppConfig(rootPath)
	if err != nil {
		panic(err)
	}

	redisClient := persist.InitRedis(appConfig.RedisConnectionString)

	_, err = redisClient.Del(context.TODO(), "waitingQueue").Result()
	if err != nil {
		panic(err)
	}

	_, err = redisClient.Del(context.TODO(), "processingQueue").Result()
	if err != nil {
		panic(err)
	}
}
