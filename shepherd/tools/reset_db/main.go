package main

import (
	"context"
	"path/filepath"

	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"github.com/el-mike/dogecrack/shepherd/internal/crack"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull"
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

	mongoClient, err := persist.InitMongo(context.TODO(), appConfig.MongoConnectionString)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := persist.GetDatabase()

	jobsCollection := db.Collection(crack.JobsCollection)

	if err := jobsCollection.Drop(context.TODO()); err != nil {
		panic(err)
	}

	instancesCollection := db.Collection(pitbull.InstancesCollection)

	if err := instancesCollection.Drop(context.TODO()); err != nil {
		panic(err)
	}
}
