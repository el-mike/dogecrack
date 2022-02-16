package main

import (
	"context"
	"path/filepath"

	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull/repositories"
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

	db := persist.GetDatabase()

	jobsCollection := db.Collection(repositories.JobsCollection)

	if err := jobsCollection.Drop(context.TODO()); err != nil {
		panic(err)
	}

	instancesCollection := db.Collection(repositories.InstancesCollection)

	if err := instancesCollection.Drop(context.TODO()); err != nil {
		panic(err)
	}
}
