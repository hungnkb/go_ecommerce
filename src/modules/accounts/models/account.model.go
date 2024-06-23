package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CredentialProvider int

const (
	PASSWORD CredentialProvider = 1
	GOOGLE   CredentialProvider = 2
	FACEBOOK CredentialProvider = 3
)

type Credential struct {
	Provider CredentialProvider `json:"provider"`
	Password string             `json:"password"`
	Account  Account            `json:"account"`
}

type Account struct {
	ID           primitive.ObjectID `json:"_id,omitempty"`
	Username     string             `json:"username" binding:"required"`
	Name         string             `json:"name" binding:"required"`
	Age          int                `json:"age" binding:"required"`
	Gender       string             `json:"gender" binding:"required"`
	Email        string             `json:"email" binding:"required"`
	Phone        string             `json:"phone" binding:"required"`
	ProvinceCode string             `json:"province_code"`
	DistrictCode string             `json:"district_code"`
	WardCode     string             `json:"ward_code"`
	Address      string             `json:"address"`
	Password     string             `json:"password" binding:"required"`
	Credentials  []Credential       `json:"credentials"`
}
