package main

import (
	"context"
	"path/filepath"
	"time"

	"github.com/el-mike/dogecrack/shepherd/config"
	"github.com/el-mike/dogecrack/shepherd/persist"
	"github.com/el-mike/dogecrack/shepherd/pitbull"
	"github.com/el-mike/dogecrack/shepherd/server"
	"github.com/el-mike/dogecrack/shepherd/vast"
)

func main() {
	appConfig, err := config.NewAppConfig()
	if err != nil {
		panic(err)
	}

	rootPath, err := filepath.Abs(".")
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

	vastManager := vast.NewVastManager(appConfig.VastApiSecret, appConfig.PitbullImage, appConfig.SSHUser, appConfig.SSHPassword, appConfig.SSHDirPath, rootPath)

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

	pitbullManager := pitbull.NewPitbullManager(vastManager)

	s := server.NewServer(pitbullManager)
	s.Run()
}
