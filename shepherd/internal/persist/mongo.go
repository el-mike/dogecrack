package persist

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const DatabaseName = "shepherd"

var mongoClient *mongo.Client

// GetMongoClient - singleton implementation for mongoClient.
func GetMongoClient() *mongo.Client {
	return mongoClient
}

// GetDatabase - helper function for getting Shepherd's database.
func GetDatabase() *mongo.Database {
	if mongoClient != nil {
		return mongoClient.Database(DatabaseName)
	}

	return nil
}

// InitMongo - inits MongoDB instance connection.
func InitMongo(ctx context.Context, connectionString string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	mongoClient = client

	return client, nil
}
