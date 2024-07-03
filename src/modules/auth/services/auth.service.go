package authService

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	httpMessage "github.com/hungnkb/go_ecommerce/src/common/httpCommon/http-error-message"
	httpStatusCode "github.com/hungnkb/go_ecommerce/src/common/httpCommon/http-status"
	Config "github.com/hungnkb/go_ecommerce/src/config"
	"github.com/hungnkb/go_ecommerce/src/modules/accountStorage"
	"github.com/hungnkb/go_ecommerce/src/modules/accounts/accountModel"
	authDto "github.com/hungnkb/go_ecommerce/src/modules/auth/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Register(db *mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var body accountModel.Account
		if error := c.ShouldBind(&body); error != nil {
			c.JSON(400, gin.H{
				"message": error.Error(),
			})
			return
		}

		data := accountStorage.InsertAccount(db, body)
		if data.Data != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"data":   data.Data,
			})
		} else {
			c.JSON(data.HttpStatusCode, gin.H{
				"status":  data.HttpStatusCode,
				"message": data.Error,
			})
		}
	}
	return fn
}

func LoginByPassword(db *mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var body authDto.Login
		if error := c.ShouldBind(&body); error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": error.Error(),
			})
		}

		data := accountStorage.GetAccountBy(db, bson.D{{Key: "username", Value: body.Username}})
		if data.Email == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  httpStatusCode.UNAUTHORIZED,
				"message": httpMessage.ERROR_ACCOUNT_NOT_FOUND,
			})
			return
		}
		var hasedPassword []byte
		for i := range data.Credentials {
			if int(data.Credentials[i].Provider) == int(accountModel.PASSWORD) {
				hasedPassword = []byte(data.Credentials[i].Password)
				break
			}
		}
		if err := verifyPassword([]byte(hasedPassword), []byte(body.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": httpMessage.ERROR_ACCOUNT_WRONG_PASSWORD,
			})
			return
		}
		accessToken, atError := AccessTokenGenerator(data)
		if atError == nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"data": gin.H{
					"accessToken": accessToken,
				},
			})
		} else {
			fmt.Println(atError)
		}
	}
	return fn
}

func AccessTokenVerify(tokenString string) (jwt.MapClaims, error) {
	secretKey := Config.Get().SecretKey
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func AccessTokenGenerator(data accountModel.Account) (string, error) {
	aWeekHour := 24 * 7
	secretKey := Config.Get().SecretKey
	fmt.Println(data)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": data.ID,
		"iss": "Ecommerce",
		"exp": time.Now().Add(time.Hour * time.Duration(aWeekHour)).UnixMilli(),
		"iat": time.Now().UnixMilli(),
	})
	token, error := claims.SignedString([]byte(secretKey))
	return token, error
}

func verifyPassword(hasedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hasedPassword, password)
}