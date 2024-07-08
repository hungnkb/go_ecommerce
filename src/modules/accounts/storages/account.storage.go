package accountStorage

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	httpMessage "github.com/hungnkb/go_ecommerce/src/common/httpCommon/http-error-message"
	"github.com/hungnkb/go_ecommerce/src/common/query"
	responseType "github.com/hungnkb/go_ecommerce/src/common/types"
	Config "github.com/hungnkb/go_ecommerce/src/config"
	accountModel "github.com/hungnkb/go_ecommerce/src/modules/accounts/models"
	"github.com/hungnkb/go_ecommerce/src/modules/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var accountCollection = "accounts"
var permissionCollection = "permissions"

func GetAccountList(client *mongo.Client, filter bson.M, page int64, limit int64, search string) responseType.StorageReponseType {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 100
	}
	skip := limit * (page - 1)
	var searchParsed bson.M = bson.M{}
	if search != "" {
		searchParsed = query.SearchParser(search)
	}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{
		Key:   "from",
		Value: "permissions",
	}, {
		Key:   "localField",
		Value: "permission_ids",
	}, {
		Key:   "foreignField",
		Value: "_id",
	}, {
		Key:   "as",
		Value: "permissions",
	}},
	}}
	skipStage := bson.D{{Key: "$skip", Value: skip}}
	limitStage := bson.D{{Key: "$limit", Value: limit}}
	searchStage := bson.D{{Key: "$match", Value: searchParsed}}
	cursor, error := client.Database(storage.DB).Collection(accountCollection).Aggregate(context.TODO(), mongo.Pipeline{lookupStage, skipStage, limitStage, searchStage})
	if error != nil {
		println(error.Error())
		return responseType.StorageReponseType{
			Data:           nil,
			Error:          error.Error(),
			HttpStatusCode: http.StatusBadRequest,
		}
	}
	var result []accountModel.Account
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

func GetAccountBy(client *mongo.Client, stages []bson.D) []accountModel.Account {
	var result []accountModel.Account
	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "permissions"},
			{Key: "localField", Value: "permission_ids"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "permissions"}},
		},
	}
	stages = append(stages, lookupStage)
	if cur, errCur := storage.GetColection(client, accountCollection).Aggregate(context.TODO(), mongo.Pipeline(stages)); errCur == nil {
		fmt.Println(123)
		cur.All(context.TODO(), &result)
	} else {
		fmt.Println(errCur.Error())
	}
	return result
}

func InsertAccount(client *mongo.Client, account accountModel.Account) responseType.StorageReponseType {
	var result accountModel.Account
	findQuery := bson.D{primitive.E{Key: "$or",
		Value: []interface{}{
			bson.D{{Key: "username", Value: account.Username}},
			bson.D{{Key: "email", Value: account.Email}},
			bson.D{{Key: "phone", Value: account.Phone}},
		},
	}}
	checkExistError := storage.GetColection(client, accountCollection).FindOne(context.TODO(), findQuery).Decode(&result)
	if checkExistError == nil {
		fmt.Println("InsertAccount/checkExistError", result)
		return responseType.StorageReponseType{
			Error:          string(httpMessage.ERROR_ACCOUNT_EXIST),
			HttpStatusCode: http.StatusConflict,
		}
	}

	hashCost := Config.Get().Salt
	hashPassword, hashError := bcrypt.GenerateFromPassword([]byte(account.Password), hashCost)
	if hashError != nil {
		fmt.Println("InsertAccount/hashError", hashError.Error())
		return responseType.StorageReponseType{
			Error:          "Hash password failed",
			HttpStatusCode: http.StatusOK,
		}
	}
	credential := &accountModel.Credential{
		Provider: accountModel.PASSWORD,
		Password: string(hashPassword),
	}

	var userPermission accountModel.Permission
	if len(account.PermissionIDs) == 0 {
		storage.GetColection(client, permissionCollection).FindOne(context.TODO(), bson.D{{Key: "key", Value: "user"}}).Decode(&userPermission)
		if objectId, err := primitive.ObjectIDFromHex(userPermission.ID); err == nil {
			account.PermissionIDs = append(account.PermissionIDs, objectId)
		}
	}

	account.Password = ""
	account.Credentials = []accountModel.Credential{*credential}
	account.CreatedAt = primitive.DateTime(time.Now().UnixMilli())
	account.UpdatedAt = primitive.DateTime(time.Now().UnixMilli())
	insertResult, insertError := storage.GetColection(client, accountCollection).InsertOne(context.TODO(), &account)
	if insertError != nil {
		fmt.Println("InsertAccount/insertError", insertError.Error())
		return responseType.StorageReponseType{
			Error:          "Something went wrong",
			HttpStatusCode: http.StatusBadRequest,
		}
	}
	account.ID = insertResult.InsertedID.(primitive.ObjectID)
	account.Permissions = append(account.Permissions, userPermission)
	return responseType.StorageReponseType{
		Data: account,
	}
}

func InsertPermissionBulk(client *mongo.Client, input []accountModel.Permission) responseType.StorageReponseType {
	var listKey []string
	for i := range input {
		if !slices.Contains(listKey, input[i].Key) {
			listKey = append(listKey, input[i].Key)
		}
	}
	check := storage.GetColection(client, permissionCollection).FindOne(context.TODO(), bson.D{{
		Key: "key", Value: bson.D{
			{Key: "$in", Value: listKey},
		},
	}})
	if check.Err() == nil {
		return responseType.StorageReponseType{
			HttpStatusCode: http.StatusConflict,
			Error:          httpMessage.ERROR_PERMISSION_KEY_EXIST,
		}
	}

	var data []accountModel.Permission
	formatInput := make([]interface{}, 0)
	for _, v := range input {
		formatInput = append(formatInput, v)
	}
	res, error := storage.GetColection(client, permissionCollection).InsertMany(context.TODO(), formatInput)
	if error != nil {
		return responseType.StorageReponseType{
			HttpStatusCode: http.StatusBadRequest,
			Error:          error.Error(),
		}
	}
	if cur, err := storage.GetColection(client, permissionCollection).Find(context.TODO(), bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: res.InsertedIDs}}}}); err == nil {
		cur.All(context.TODO(), &data)
	}
	fmt.Println(data)
	return responseType.StorageReponseType{
		Data: data,
	}
}
