package persist

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client

// GetMongoClient - singleton implementation for mongoClient.
func GetMongoClient() *mongo.Client {
	return mongoClient
}

// InitMongo - inits MongoDB instance connection.
func InitMongo(ctx context.Context, username, password, host, port string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(getURI(username, password, host, port))

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

func getURI(username, password, host, port string) string {
	return "mongodb://" + username + ":" + password + "@" + host + ":" + port
}
