package main

import (
	"context"
	"path/filepath"

	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"github.com/el-mike/dogecrack/shepherd/internal/server"
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

	sshIp, err := vast.GetFakeVastIp(rootPath, 1)
	if err != nil {
		panic(err)
	}

	if err := vast.AddSSHFingerprint(rootPath, sshIp, appConfig.SSHDirPath); err != nil {
		panic(err)
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

	s := server.NewServer()
	s.Run()
}
