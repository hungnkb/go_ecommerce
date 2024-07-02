package authDto

import "github.com/hungnkb/go_ecommerce/src/modules/accounts/models"

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginReponse struct {
	AccessToken  string         `json:"accessToken"`
	RefreshToken string         `json:"refreshToken"`
	Account      models.Account `json:"account"`
}
