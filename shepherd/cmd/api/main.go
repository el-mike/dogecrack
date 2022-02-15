package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"github.com/el-mike/dogecrack/shepherd/internal/core"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"github.com/el-mike/dogecrack/shepherd/internal/vast"
)

func main() {
	// This assumes that we run the application from the project's root directory,
	// NOT /cmd/api.
	// This approach helps with running the app in Docker containers, where built app
	// is no longer in /cmd/api directory.
	rootPath, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}

	handleFakeVast := false

	for _, arg := range os.Args[1:] {
		if arg == "handleFakeVast" {
			handleFakeVast = true
		}
	}

	appConfig, err := config.NewAppConfig(rootPath)
	if err != nil {
		panic(err)
	}

	if handleFakeVast {
		setupFakeVast(appConfig)
	}

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

	server := core.NewServer(appConfig.APIPort, appConfig.OriginAllowed)
	server.Run()
}

func setupFakeVast(config *config.AppConfig) {
	sshIp, err := vast.GetFakeVastIp(config.RootPath, 1)
	if err != nil {
		panic(err)
	}

	if err := vast.AddSSHFingerprint(config.RootPath, sshIp, config.SSHDirPath); err != nil {
		panic(err)
	}
}
