package storage

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/hungnkb/go_ecommerce/src/modules/accounts/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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

func InsertAccount(client *mongo.Client, account models.Account) interface{} {
	var result models.Account
	findQuery := bson.D{primitive.E{Key: "$or",
		Value: []interface{}{
			bson.D{{Key: "username", Value: account.Username}},
			bson.D{{Key: "email", Value: account.Email}},
			bson.D{{Key: "phone", Value: account.Phone}},
		},
	}}
	checkExistError := getColection(client, DB).FindOne(context.TODO(), findQuery).Decode(&result)
	if checkExistError == nil {
		fmt.Println(result)
		return nil
	}
	fmt.Println(123123)
	hashCost, _ := strconv.Atoi(os.Getenv("salt"))
	hashPassword, hashError := bcrypt.GenerateFromPassword([]byte(account.Password), hashCost)
	if hashError != nil {
		return nil
	}
	credential := &models.Credential{
		Provider: models.PASSWORD,
		Password: string(hashPassword),
	}
	account.Password = ""
	account.Credentials = []models.Credential{*credential}
	fmt.Println(123123, account)
	insertResult, insertError := getColection(client, DB).InsertOne(context.TODO(), &account)
	if insertError != nil {
		return nil
	}
	account.ID = insertResult.InsertedID.(primitive.ObjectID)
	return account
}
