package documentModel

import "go.mongodb.org/mongo-driver/bson/primitive"

type Document struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Type string             `json:"type"`
	Name string             `json:"string"`
	Url  string             `json:"url"`
}
