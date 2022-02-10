package main

import (
	"context"
	"path/filepath"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull"
	"github.com/el-mike/dogecrack/shepherd/internal/vast"
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

	vastManager := vast.NewVastManager(appConfig.VastApiSecret, appConfig.PitbullImage, appConfig.SSHUser, appConfig.SSHPassword, appConfig.SSHDirPath, rootPath)

	mongoClient, err := persist.InitMongo(context.TODO(), appConfig.MongoUser, appConfig.MongoPassword, appConfig.MongoHost, appConfig.MongoPort)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	persist.InitRedis(appConfig.RedisHost, appConfig.RedisPort)

	manager := pitbull.NewManager(vastManager)
	runner := pitbull.NewRunner(manager)

	collector := pitbull.NewCollector(15 * time.Second)
	dispatcher := pitbull.NewJobDispatcher(runner, 15*time.Second)

	go collector.Start()

	dispatcher.Start()
}
