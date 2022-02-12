package main

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull"
	"github.com/el-mike/dogecrack/shepherd/internal/vast"
)

func main() {
	logger := common.NewLogger("Startup", os.Stdout, os.Stderr)

	rootPath, err := filepath.Abs("../../")
	if err != nil {
		logger.Err.Println(err)
		panic(err)
	}

	appConfig, err := config.NewAppConfig(rootPath)
	if err != nil {
		logger.Err.Println(err)
		panic(err)
	}

	vastManager := vast.NewVastManager(appConfig.VastApiSecret, appConfig.PitbullImage, appConfig.SSHUser, appConfig.SSHPassword, appConfig.SSHDirPath, rootPath)

	mongoClient, err := persist.InitMongo(context.TODO(), appConfig.MongoUser, appConfig.MongoPassword, appConfig.MongoHost, appConfig.MongoPort)
	if err != nil {
		logger.Err.Println(err)
		panic(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			logger.Err.Println(err)
			panic(err)
		}
	}()

	persist.InitRedis(appConfig.RedisHost, appConfig.RedisPort)

	instanceManager := pitbull.NewInstanceManager(vastManager)
	jobManager := pitbull.NewJobManager(instanceManager)
	runner := pitbull.NewJobRunner(instanceManager)

	// On service start, we want to reschedule all jobs from "processingQueue",
	// as since the worker has been restarted, there is no thread working on those tasks anymore.
	jobsIds, err := jobManager.RescheduleProcessingJobs()
	if err != nil {
		logger.Err.Printf("processing jobs rescheduling failed. reason: %v\n", err)
	}

	if len(jobsIds) > 0 {
		logger.Info.Printf("%d jobs have been rescheduled.", len(jobsIds))
	}

	collector := pitbull.NewInstanceCollector(instanceManager, 15*time.Second)
	dispatcher := pitbull.NewJobDispatcher(instanceManager, runner, 15*time.Second)

	go collector.Start()

	dispatcher.Start()
}
