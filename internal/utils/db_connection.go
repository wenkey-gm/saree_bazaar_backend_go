package utils

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func DbConnection() *mongo.Client {
	MongoUrl := os.Getenv("MONGO_URL")
	if MongoUrl == "" {
		MongoUrl = MONGO_URL
	}
	clientOptions := options.Client().ApplyURI(MongoUrl) // mongodb://localhost:27017

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

func ConnectMongoDbCollection(client *mongo.Client, dbName string, collectionName string) *mongo.Collection {

	return client.Database(dbName).Collection(collectionName)
}
