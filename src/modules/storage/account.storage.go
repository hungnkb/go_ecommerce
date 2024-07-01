package storage

import (
	"context"
	"fmt"
	"os"
	"strconv"

	httpMessage "github.com/hungnkb/go_ecommerce/src/common/httpCommon/http-error-message"
	httpStatusCode "github.com/hungnkb/go_ecommerce/src/common/httpCommon/http-status"
	responseType "github.com/hungnkb/go_ecommerce/src/common/types"
	"github.com/hungnkb/go_ecommerce/src/modules/accounts/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var DB = os.Getenv("DB_NAME")
var accountCollection = "accounts"

func GetAccountList(client *mongo.Client) []models.Account {
	cursor, error := client.Database(DB).Collection(accountCollection).Find(context.TODO(), models.Account{})
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

func InsertAccount(client *mongo.Client, account models.Account) responseType.StorageReponseType {
	var result models.Account
	findQuery := bson.D{primitive.E{Key: "$or",
		Value: []interface{}{
			bson.D{{Key: "username", Value: account.Username}},
			bson.D{{Key: "email", Value: account.Email}},
			bson.D{{Key: "phone", Value: account.Phone}},
		},
	}}
	checkExistError := getColection(client, accountCollection).FindOne(context.TODO(), findQuery).Decode(&result)
	if checkExistError == nil {
		fmt.Println("InsertAccount/checkExistError", result)
		return responseType.StorageReponseType{
			Error:          string(httpMessage.ERROR_ACCOUNT_EXIST),
			HttpStatusCode: int(httpStatusCode.CONFLICT),
		}
	}
	hashCost, _ := strconv.Atoi(os.Getenv("salt"))
	hashPassword, hashError := bcrypt.GenerateFromPassword([]byte(account.Password), hashCost)
	if hashError != nil {
		fmt.Println("InsertAccount/hashError", hashError.Error())
		return responseType.StorageReponseType{
			Error:          "Hash password failed",
			HttpStatusCode: int(httpStatusCode.OK),
		}
	}
	credential := &models.Credential{
		Provider: models.PASSWORD,
		Password: string(hashPassword),
	}
	account.Password = ""
	account.Credentials = []models.Credential{*credential}
	insertResult, insertError := getColection(client, accountCollection).InsertOne(context.TODO(), &account)
	if insertError != nil {
		fmt.Println("InsertAccount/insertError", insertError.Error())
		return responseType.StorageReponseType{
			Error:          "Something went wrong",
			HttpStatusCode: 400,
		}
	}
	account.ID = insertResult.InsertedID.(primitive.ObjectID)
	return responseType.StorageReponseType{
		Data: account,
	}
}
