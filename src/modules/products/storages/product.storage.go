package productStorage

import (
	"context"
	"fmt"
	"net/http"

	httpMessage "github.com/hungnkb/go_ecommerce/src/common/httpCommon/http-error-message"
	responseType "github.com/hungnkb/go_ecommerce/src/common/types"
	accountModel "github.com/hungnkb/go_ecommerce/src/modules/accounts/models"
	productModel "github.com/hungnkb/go_ecommerce/src/modules/products/models"

	"github.com/hungnkb/go_ecommerce/src/modules/storage"
	"github.com/hungnkb/go_ecommerce/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	productCollectionName         = "products"
	productMetadataCollectionName = "product_metadata"
	attributeCollectionName       = "attributes"
)

type ProductInput struct {
	Product         productModel.Product           `json:"product"`
	ProductMetadata []productModel.ProductMetadata `json:"productMetadata"`
}

func InsertProduct(db *mongo.Client, input ProductInput, account accountModel.Account) responseType.StorageReponseType {
	slug, idProduct := utils.SlugGeneratorWithId(input.Product.Name, primitive.NilObjectID)
	input.Product.ID = idProduct
	input.Product.Slug = slug
	input.Product.AccountId = account.ID
	_, productErr := storage.GetColection(db, productCollectionName).InsertOne(context.TODO(), input.Product)
	if productErr != nil {
		return responseType.StorageReponseType{
			Error:          httpMessage.ERROR_INSERT_PRODUCT,
			HttpStatusCode: http.StatusBadRequest,
		}
	}
	if len(input.ProductMetadata) > 0 {
		var productMetadataFormated []interface{}
		for _, metadata := range input.ProductMetadata {
			metadata.ProductId = idProduct
			productMetadataFormated = append(productMetadataFormated, metadata)
		}
		resultMetadata, errMetadata := storage.GetColection(db, productMetadataCollectionName).InsertMany(context.TODO(), productMetadataFormated)
		if errMetadata == nil {
			fmt.Print(666, errMetadata)
			return responseType.StorageReponseType{
				Error:          httpMessage.ERROR_INSERT_PRODUCT_METADATA,
				HttpStatusCode: http.StatusBadRequest,
			}
		}
		fmt.Println(777, resultMetadata.InsertedIDs)
	}

	return responseType.StorageReponseType{}
}

// func InsertBulkProductMetadata(db *mongo.Client, input []ProductMetadata) {

// }

func InsertAttributeBulk(db *mongo.Client, input []productModel.ProductAttribute, account accountModel.Account) responseType.StorageReponseType {
	var inputFormat []interface{}
	for _, item := range input {
		if account.ID != primitive.NilObjectID {
			item.AccountId = account.ID
		}
		inputFormat = append(inputFormat, item)
	}
	res, err := storage.GetColection(db, attributeCollectionName).InsertMany(context.TODO(), inputFormat)
	if err == nil {
		return responseType.StorageReponseType{
			Error:          httpMessage.ERROR_INSERT_PRODUCT_ATTRIBUTE,
			HttpStatusCode: http.StatusBadRequest,
		}
	}
	var attributes []productModel.ProductAttribute
	filter := bson.D{{Key: "$_id", Value: bson.D{{Key: "_id", Value: res.InsertedIDs}}}}
	cur, errCur := storage.GetColection(db, attributeCollectionName).Find(context.TODO(), filter)
	if errCur != nil {
		cur.All(context.TODO(), attributes)
	}
	return responseType.StorageReponseType{
		Data:           attributes,
		Error:          "",
		HttpStatusCode: http.StatusOK,
	}
}
