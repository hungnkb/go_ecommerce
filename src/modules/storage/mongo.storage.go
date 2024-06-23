package storage

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func NewMongoStorage() *mongo.Client {
	mongodbUrl := os.Getenv("MONGODB_URL")
	var ctx = context.TODO()
	options2 := options.Client().ApplyURI(mongodbUrl)
	client, err := mongo.Connect(ctx, options2)
	if err != nil {
		println("MongoDb connect failed!")
		return nil
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	println("MongoDb connected!")
	return client
}

func getColection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(DB).Collection(collectionName)
	return collection
}
