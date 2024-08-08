package accountModel

type Permission struct {
	ID   string `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" binding:"required"`
	Key  string `json:"key" binding:"required"`
}
