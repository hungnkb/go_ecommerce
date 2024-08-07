package accountModel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CredentialProvider int

const (
	PASSWORD CredentialProvider = 1
	GOOGLE   CredentialProvider = 2
	FACEBOOK CredentialProvider = 3
)

type Credential struct {
	Provider CredentialProvider `json:"provider"`
	Password string             `json:"password"`
}

type Account struct {
	ID            primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Username      string               `json:"username" binding:"required"`
	Name          string               `json:"name"`
	Dob           primitive.DateTime   `json:"dob"`
	Gender        string               `json:"gender"`
	Email         string               `json:"email"`
	Phone         string               `json:"phone"`
	ProvinceCode  string               `json:"provinceCode" bson:"province_code"`
	DistrictCode  string               `json:"districtCode" bson:"district_code"`
	WardCode      string               `json:"wardCode" bson:"ward_code"`
	Address       string               `json:"address"`
	Password      string               `json:"password,omitempty" binding:"required"`
	IsShop        bool                 `json:"is_shop"`
	CreatedAt     primitive.DateTime   `json:"createdAt" bson:"created_at"`
	UpdatedAt     primitive.DateTime   `json:"updatedAt" bson:"updated_at"`
	Credentials   []Credential         `json:"-"`
	PermissionIDs []primitive.ObjectID `json:"permissionIds,omitempty" bson:"permission_ids,omitempty"`
	Permissions   []Permission         `json:"permissions"`
}
