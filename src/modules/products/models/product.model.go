package productModel

import (
	documentModel "github.com/hungnkb/go_ecommerce/src/modules/documents/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductAttribute struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name"`
	Key       string             `json:"key"`
	AccountId primitive.ObjectID `json:"accountId,omitempty" bson:"account_id,omitempty"`
}

type ProductMetadata struct {
	ID          primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Sku         string                 `json:"sku"`
	AttributeId primitive.ObjectID     `json:"attributeId" bson:"attribute_id"`
	Atribute    ProductAttribute       `json:"-"`
	ProductId   primitive.ObjectID     `json:"productId" bson:"product_id"`
	Product     Product                `json:"-"`
	Value       string                 `json:"value"`
	DocumentId  primitive.ObjectID     `json:"documentId,omitempty" bson:"document_id,omitempty"`
	Document    documentModel.Document `json:"-"`
	IsThumbnail bool                   `json:"isThumbnail" bson:"is_thumbnail"`
	Quantity    int                    `json:"quantity,omitempty" bson:"-" binding:"required"`
}

type Product struct {
	ID              primitive.ObjectID       `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string                   `json:"name" binding:"required"`
	Slug            string                   `json:"slug"`
	Price           float64                  `json:"price" binding:"required"`
	Description     string                   `json:"description"`
	AccountId       primitive.ObjectID       `json:"accountId" bson:"account_id"`
	ProductMetadata []ProductMetadata        `json:"productMetadata"`
	DocumentIds     []primitive.ObjectID     `json:"documentIds" bson:"document_ids"`
	Documents       []documentModel.Document `json:"documents"`
	ThumbnailId     primitive.ObjectID       `json:"thumbnailId" bson:"thumbnail_id"`
}
