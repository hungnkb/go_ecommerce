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
	categoryCollectionName        = "categories"
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
		fmt.Println(productErr.Error())
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
		_, errMetadata := storage.GetColection(db, productMetadataCollectionName).InsertMany(context.TODO(), productMetadataFormated)
		if errMetadata != nil {
			return responseType.StorageReponseType{
				Error:          httpMessage.ERROR_INSERT_PRODUCT_METADATA,
				HttpStatusCode: http.StatusBadRequest,
			}
		}
	}
	var data productModel.Product
	filter := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: idProduct}}}}
	lookup := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "product_metadata"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "product_id"},
			{Key: "as", Value: "productMetadata"},
		}},
	}
	cursor, errGetResult := storage.GetColection(db, productCollectionName).Aggregate(context.TODO(), mongo.Pipeline{filter, lookup})
	if errGetResult != nil {
		fmt.Println(errGetResult.Error())
		return responseType.StorageReponseType{
			HttpStatusCode: http.StatusBadRequest,
			Error:          errGetResult.Error(),
		}
	}
	if cursor.Next(context.TODO()) {
		err := cursor.Decode(&data)
		if err != nil {
			fmt.Println(err.Error())
			return responseType.StorageReponseType{
				HttpStatusCode: http.StatusBadRequest,
				Error:          err.Error(),
			}
		}
	}

	return responseType.StorageReponseType{
		HttpStatusCode: http.StatusOK,
		Data:           data,
	}
}

func ProductGetList(db *mongo.Client, page, limit int64, keywords string) responseType.StorageReponseType {
	var result []productModel.Product
	lookupDocumentStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "documents"},
			{Key: "localField", Value: "document_ids"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "documents"},
		}},
	}
	lookupMedataStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "product_metadata"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "product_id"},
			{Key: "as", Value: "productMetadata"},
		}},
	}
	limitStage := bson.D{
		{Key: "$limit", Value: limit},
	}
	skipStage := bson.D{
		{Key: "$skip", Value: (page - 1) * limit},
	}
	cursor, err := storage.GetColection(db, productCollectionName).Aggregate(context.TODO(), mongo.Pipeline{lookupDocumentStage, lookupMedataStage, limitStage, skipStage})
	if err != nil {
		fmt.Println(err.Error())
	}
	cursor.All(context.TODO(), &result)
	return responseType.StorageReponseType{
		HttpStatusCode: http.StatusOK,
		Data:           result,
	}
}

// func InsertBulkProductMetadata(db *mongo.Client, input []ProductMetadata) {

// }

func InsertAttributeBulk(db *mongo.Client, input []productModel.ProductAttribute, account accountModel.Account) responseType.StorageReponseType {
	var inputFormat []interface{}
	for _, item := range input {
		if account.ID != primitive.NilObjectID && account.IsShop {
			item.AccountId = account.ID
		}
		inputFormat = append(inputFormat, item)
	}
	res, err := storage.GetColection(db, attributeCollectionName).InsertMany(context.TODO(), inputFormat)
	if err != nil {
		return responseType.StorageReponseType{
			Error:          httpMessage.ERROR_INSERT_PRODUCT_ATTRIBUTE,
			HttpStatusCode: http.StatusBadRequest,
		}
	}
	var attributes []productModel.ProductAttribute
	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: res.InsertedIDs}}}}
	cur, errFind := storage.GetColection(db, attributeCollectionName).Find(context.TODO(), filter)
	if errFind == nil {
		cur.All(context.TODO(), &attributes)
	}
	fmt.Println(attributes)
	return responseType.StorageReponseType{
		Data:           attributes,
		Error:          "",
		HttpStatusCode: http.StatusOK,
	}
}

func GetProductBySlug(db *mongo.Client, slug string) responseType.StorageReponseType {
	var data productModel.Product
	// filter := bson.D{{Key: "$match", Value: bson.D{{Key: "slug", Value: slug}}}}
	// documentLookup := bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.D{
	// 		{Key: "from", Value: "documents"},
	// 		{Key: "localField", Value: "document_ids"},
	// 		{Key: "foreignField", Value: "_id"},
	// 		{Key: "as", Value: "documents"},
	// 	},
	// }}
	// metadataLookup := bson.D{{Key: "$lookup", Value: bson.D{
	// 	{Key: "from", Value: "product_metadata"},
	// 	{Key: "localField", Value: "_id"},
	// 	{Key: "foreignField", Value: "product_id"},
	// 	{Key: "as", Value: "productmetadata"},
	// },
	// }}
	// attributesLookup := bson.D{{Key: "$lookup", Value: bson.D{
	// 	{Key: "from", Value: "attributes"},
	// 	{Key: "localField", Value: "product_metadata.attribute_id"},
	// 	{Key: "foreignField", Value: "_id"},
	// 	{Key: "as", Value: "product_metadata.attributes"},
	// }}}
	// unwind := bson.D{{Key: "$unwind", Value: "$product_metadata"}}
	// unwind2 := bson.D{{Key: "$unwind", Value: "$product_metadata.attributes"}}

	pipeline := bson.A{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "documents"},
			{Key: "localField", Value: "document_ids"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "documents"},
		}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "as", Value: "product_metadata"},
					{Key: "from", Value: "product_metadata"},
					{Key: "foreignField", Value: "product_id"},
					{Key: "localField", Value: "_id"},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$product_metadata"}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "as", Value: "product_metadata.attributes"},
					{Key: "from", Value: "attributes"},
					{Key: "foreignField", Value: "_id"},
					{Key: "localField", Value: "product_metadata.attribute_id"},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$product_metadata.attributes"}}}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "slug", Value: slug}}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "productmetadata", Value: bson.D{{Key: "$push", Value: "$product_metadata"}}},
			{Key: "allFields", Value: bson.D{{Key: "$first", Value: "$$ROOT"}}},
		}}},
		bson.D{{Key: "$replaceRoot", Value: bson.D{
			{Key: "newRoot", Value: bson.D{{Key: "$mergeObjects", Value: bson.A{"$allFields", bson.D{{Key: "productMetadata", Value: "$productmetadata"}}}}}},
		}}},
	}

	cursor, err := storage.GetColection(db, productCollectionName).Aggregate(context.TODO(), pipeline)
	if err != nil {
		fmt.Println(err.Error())
		return responseType.StorageReponseType{
			HttpStatusCode: http.StatusBadRequest,
			Error:          err.Error(),
		}
	}
	cursor.Next(context.TODO())
	cursor.Decode(&data)
	return responseType.StorageReponseType{
		Data:           data,
		Error:          "",
		HttpStatusCode: http.StatusOK,
	}
}
