package main

import (
	"context"
	"path/filepath"
	"time"

	"github.com/el-mike/dogecrack/shepherd/config"
	"github.com/el-mike/dogecrack/shepherd/persist"
	"github.com/el-mike/dogecrack/shepherd/server"
	"github.com/el-mike/dogecrack/shepherd/vast"
)

func main() {
	appConfig, err := config.NewAppConfig()
	if err != nil {
		panic(err)
	}

	path, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	sshIp := vast.GetFakeVastIp(path)

	client, err := vast.NewVastClient(appConfig.SSHUser, appConfig.SSHPassword, appConfig.SSHDirPath, sshIp)
	if err != nil {
		panic(err)
	}

	err = client.Connect()
	if err != nil {
		panic(err)
	}

	defer client.Close()

	manager := vast.NewVastManager(appConfig.VastApiSecret)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := persist.InitMongo(ctx, appConfig.MongoUser, appConfig.MongoPassword, appConfig.MongoHost, appConfig.MongoPort)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	s := server.NewServer(manager, client)
	s.Run()
}
