package productModel

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" binding:"required"`
	ParentID primitive.ObjectID `json:"parentId,omitempty" bson:"parent_id,omitempty"`
	Slug     string             `json:"slug"`
	Img      string             `json:"img"`
}
