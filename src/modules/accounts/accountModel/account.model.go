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
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username     string             `json:"username" binding:"required"`
	Name         string             `json:"name" binding:"required"`
	Dob          primitive.DateTime `json:"dob" binding:"required"`
	Gender       string             `json:"gender" binding:"required"`
	Email        string             `json:"email" binding:"required"`
	Phone        string             `json:"phone" binding:"required"`
	ProvinceCode string             `json:"provinceCode" bson:"province_code"`
	DistrictCode string             `json:"districtCode" bson:"district_code"`
	WardCode     string             `json:"wardCode" bson:"ward_code"`
	Address      string             `json:"address"`
	Password     string             `json:"password,omitempty" binding:"required"`
	CreatedAt    primitive.DateTime `json:"createdAt" bson:"created_at"`
	UpdatedAt    primitive.DateTime `json:"updatedAt" bson:"updated_at"`
	Credentials  []Credential       `json:"credentials"`
	PermissionID primitive.ObjectID `json:"permissionId,omitempty" bson:"permission_id,omitempty"`
	Permissions  Permission         `json:"permissions"`
}