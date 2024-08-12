package productStorage

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	httpStatusCode "github.com/hungnkb/go_ecommerce/src/common/httpCommon/http-status"
	responseType "github.com/hungnkb/go_ecommerce/src/common/types"
	productModel "github.com/hungnkb/go_ecommerce/src/modules/products/models"
	"github.com/hungnkb/go_ecommerce/src/modules/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertCategory(db *mongo.Client, input productModel.Category) responseType.StorageReponseType {
	input.Slug = strings.ReplaceAll(strings.ToLower(input.Name), " ", "-") + "-" + strconv.Itoa(int(time.Now().UnixMilli()))
	result, err := storage.GetColection(db, categoryCollectionName).InsertOne(context.TODO(), input)
	if err != nil {
		return responseType.StorageReponseType{
			Error:          err.Error(),
			HttpStatusCode: int(httpStatusCode.BAD_REQUEST),
		}
	}
	categoryID := result.InsertedID.(primitive.ObjectID)
	newCategory := FindCategoryById(db, categoryID)
	return responseType.StorageReponseType{
		HttpStatusCode: int(httpStatusCode.OK),
		Data:           newCategory,
	}
}

func FindCategoryById(db *mongo.Client, ID primitive.ObjectID) productModel.Category {
	var result productModel.Category
	cursor := storage.GetColection(db, categoryCollectionName).FindOne(context.TODO(), bson.D{{Key: "_id", Value: ID}})
	err := cursor.Decode(&result)
	if err != nil {
		fmt.Println(err.Error())
	}
	return result
}

func GetListCategory(db *mongo.Client, page, limit int, parentId string) responseType.StorageReponseType {
	var data []productModel.Category
	limitStage := bson.D{{Key: "$limit", Value: limit}}
	skipStage := bson.D{{Key: "$skip", Value: (page - 1) * limit}}
	var filterStage primitive.D
	if parentId == "null" {
		filterStage = bson.D{{Key: "$match", Value: bson.D{{Key: "parent_id", Value: nil}}}}
	} else {
		parentObjId, _ := primitive.ObjectIDFromHex(parentId)
		filterStage = bson.D{{Key: "$match", Value: bson.D{{Key: "parent_id", Value: parentObjId}}}}
	}
	cursor, err := storage.GetColection(db, categoryCollectionName).Aggregate(context.TODO(), mongo.Pipeline{limitStage, skipStage, filterStage})
	if err != nil {
		return responseType.StorageReponseType{
			Error:          err.Error(),
			HttpStatusCode: int(httpStatusCode.BAD_REQUEST),
		}
	}
	cursor.All(context.TODO(), &data)
	return responseType.StorageReponseType{
		Data:           data,
		HttpStatusCode: int(httpStatusCode.OK),
	}
}
