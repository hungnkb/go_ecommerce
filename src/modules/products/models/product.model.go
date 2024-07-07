package productModel

import (
	documentModel "github.com/hungnkb/go_ecommerce/src/modules/documents/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductAttribute struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name"`
	AccountId primitive.ObjectID `json:"accountId" bson:"account_id`
}

type ProductMetadata struct {
	ID           primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Sku          string                 `json:"sku"`
	AttributeIds primitive.ObjectID     `json:"attributeIds" bson:"attribute_ids"`
	Atributes    ProductAttribute       `json:"-"`
	ProductId    primitive.ObjectID     `json:"productId" bson:"product_id"`
	Product      Product                `json:"-"`
	Value        string                 `json:"value"`
	DocumentId   primitive.ObjectID     `json:"documentId" bson:"document_id"`
	Document     documentModel.Document `json:"-"`
	IsThumbnail  bool                   `json:"isThumbnail" bson:"is_thumbnail"`
	Quantity     int                    `json:"quantity" bson:"quantity"`
}

type Product struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string             `json:"name" binding:"required"`
	Slug            string             `json:"slug"`
	Price           float64            `json:"price"`
	AccountId       primitive.ObjectID `json:"accountId" bson:"account_id"`
	ProductMetadata []ProductMetadata  `json:"-"`
	Quantity        int                `json:"quantity,omitempty" bson:"-"`
}
