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

	mongoCtx, mongoCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoCancel()

	mongoClient, err := persist.InitMongo(mongoCtx, appConfig.MongoUser, appConfig.MongoPassword, appConfig.MongoHost, appConfig.MongoPort)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(mongoCtx); err != nil {
			panic(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	persist.InitRedis(ctx, appConfig.RedisHost, appConfig.RedisPort)

	manager := pitbull.NewPitbullManager(vastManager)
	runner := pitbull.NewPitbullRunner(manager)

	dispatcher := pitbull.NewJobDispatcher(runner, 15*time.Second)

	dispatcher.Start()
}
