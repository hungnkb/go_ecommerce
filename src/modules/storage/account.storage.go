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

var DB = os.Getenv("MONGODB_DBNAME")

const AccountCollectionName string = "accounts"

func GetAccountList(client *mongo.Client) []models.Account {
	fmt.Println(DB)
	cursor, error := client.Database(DB).Collection(AccountCollectionName).Find(context.TODO(), models.Account{})
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

func GetAccountBy(client *mongo.Client, filter *bson.M) models.Account {
	var result models.Account
	getColection(client, DB).FindOne(context.TODO(), filter).Decode(&result)
	return result
}

func InsertAccount(client *mongo.Client, account models.Account) interface{} {
	var result models.Account
	fmt.Println(account.Username)
	findQuery := bson.M{"$or": []interface{}{
		bson.M{"username": account.Username},
		bson.M{"email": account.Email},
		bson.M{"phone": account.Phone}}}

	checkExistError := getColection(client, DB).FindOne(context.TODO(), findQuery).Decode(&result)
	if checkExistError == nil {
		return nil
	}
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
	insertResult, insertError := getColection(client, AccountCollectionName).InsertOne(context.TODO(), &account)
	if insertError != nil {
		return nil
	}
	account.ID = insertResult.InsertedID.(primitive.ObjectID)
	return account
}
