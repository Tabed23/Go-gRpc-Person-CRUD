package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var (
	ctx = context.Background()
	DB *mongo.Client = ConnectDB()
)
func ConnectDB() *mongo.Client {
	clientOptions := options.Client()
	clientOptions.ApplyURI("mongodb://mongo:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	client.Database("PersonGrpc").CreateCollection(ctx, "persons")
	return client
}

