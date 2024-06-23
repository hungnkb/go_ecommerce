package storage

import (
	"context"
	"fmt"

	"github.com/hungnkb/go_ecommerce/src/modules/accounts/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const DB = "ecommerce"

func GetAccountList(client *mongo.Client) []models.Account {
	cursor, error := client.Database(DB).Collection("accounts").Find(context.TODO(), models.Account{})
	if error != nil {
		println(error.Error())
		return nil
	}
	var result []models.Account
	if errCursor := cursor.All(context.TODO(), result); errCursor == nil {
		println(errCursor)
		return nil
	}
	return result
}

func GetAccountBy(client *mongo.Client, filter *models.Account) models.Account {
	var result models.Account
	getColection(client, DB).FindOne(context.TODO(), filter).Decode(&result)
	return result
}

func InsertAccount(client *mongo.Client, account models.Account) {
	result, err := getColection(client, DB).InsertOne(context.TODO(), account)
	if err != nil {
		println(err.Error())
	}
	findQuery := bson.D{{"$or",
		bson.D{
			bson.D{{"username", account.Username}},
			bson.D{{"email", account.Email}},
		},
	}}
	checkExist := getColection(client, DB).FindOne(context.TODO(), findQuery)
	fmt.Println(result)
}
