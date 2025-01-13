package main

import (
	"context"
	"log"
	"path/filepath"

	"github.com/el-mike/dogecrack/shepherd/database/seeds"
	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
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

	mongoClient, err := persist.InitMongo(context.TODO(), appConfig.MongoUser, appConfig.MongoPassword, appConfig.MongoHost, appConfig.MongoPort)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	usersSeed := seeds.NewUsersSeed()
	appSettingsSeed := seeds.NewAppSettingsSeed()

	err = usersSeed.Execute()
	if err != nil {
		log.Fatal(err)
	}

	err = appSettingsSeed.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
