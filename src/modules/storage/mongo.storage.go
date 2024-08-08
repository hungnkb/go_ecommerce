package storage

import (
	"context"
	"github.com/fatih/color"
	Config "github.com/hungnkb/go_ecommerce/src/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var DB = Config.Get().DbName

func NewMongoStorage() *mongo.Client {
	mongodbUrl := Config.Get().MongoDbUrl
	var ctx = context.TODO()
	options2 := options.Client().ApplyURI(mongodbUrl)
	client, err := mongo.Connect(ctx, options2)
	if err != nil {
		color.Green("MongoDb connect failed!")
		return nil
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	color.Green("MongoDb connected!")
	return client
}

func GetColection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(DB).Collection(collectionName)
	return collection
}
