package storage

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	httpMessage "github.com/hungnkb/go_ecommerce/src/common/httpCommon/http-error-message"
	responseType "github.com/hungnkb/go_ecommerce/src/common/types"
	"github.com/hungnkb/go_ecommerce/src/modules/accounts/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var DB = os.Getenv("DB_NAME")
var accountCollection = "accounts"

func GetAccountList(client *mongo.Client, filter bson.D, page int64, limit int64) responseType.StorageReponseType {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 100
	}
	skip := limit * (page - 1)
	options := options.Find().SetLimit(limit).SetSkip(skip)
	cursor, error := client.Database(DB).Collection(accountCollection).Find(context.TODO(), filter, options)
	if error != nil {
		println(error.Error())
		return responseType.StorageReponseType{
			Data:  nil,
			Error: error.Error(),
		}
	}
	var result []models.Account
	if errCursor := cursor.All(context.TODO(), &result); errCursor == nil {
		return responseType.StorageReponseType{
			Data: result,
		}
	} else {
		return responseType.StorageReponseType{
			Data:  nil,
			Error: errCursor.Error(),
		}
	}

}

func GetAccountBy(client *mongo.Client, filter bson.D) models.Account {
	var result models.Account
	getColection(client, accountCollection).FindOne(context.TODO(), filter).Decode(&result)
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
			HttpStatusCode: http.StatusConflict,
		}
	}
	hashCost, _ := strconv.Atoi(os.Getenv("salt"))
	hashPassword, hashError := bcrypt.GenerateFromPassword([]byte(account.Password), hashCost)
	if hashError != nil {
		fmt.Println("InsertAccount/hashError", hashError.Error())
		return responseType.StorageReponseType{
			Error:          "Hash password failed",
			HttpStatusCode: http.StatusOK,
		}
	}
	credential := &models.Credential{
		Provider: models.PASSWORD,
		Password: string(hashPassword),
	}
	account.Password = ""
	account.Credentials = []models.Credential{*credential}
	account.CreatedAt = primitive.DateTime(time.Now().UnixMilli())
	account.UpdatedAt = primitive.DateTime(time.Now().UnixMilli())
	insertResult, insertError := getColection(client, accountCollection).InsertOne(context.TODO(), &account)
	if insertError != nil {
		fmt.Println("InsertAccount/insertError", insertError.Error())
		return responseType.StorageReponseType{
			Error:          "Something went wrong",
			HttpStatusCode: http.StatusBadRequest,
		}
	}
	account.ID = insertResult.InsertedID.(primitive.ObjectID)
	fmt.Println(insertResult.InsertedID.(primitive.ObjectID))
	fmt.Println(account)
	return responseType.StorageReponseType{
		Data: account,
	}
}
