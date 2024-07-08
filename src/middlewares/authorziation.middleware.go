package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	httpMessage "github.com/hungnkb/go_ecommerce/src/common/httpCommon/http-error-message"
	accountStorage "github.com/hungnkb/go_ecommerce/src/modules/accounts/storages"
	authService "github.com/hungnkb/go_ecommerce/src/modules/auth/handlers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthGuard(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.Request.Header["Authorization"]
		if len(authorizationHeader) < 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": httpMessage.ERROR_MISSING_INVALID_TOKEN,
			})
			c.Abort()
		}
		authorization := c.Request.Header["Authorization"][0]
		token := strings.SplitAfter(authorization, "Bearer ")
		if len(token) != 2 || token[0] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": httpMessage.ERROR_MISSING_INVALID_TOKEN,
			})
			c.Abort()
		}
		payload, verifyError := authService.AccessTokenVerify(string(token[1]))
		if verifyError != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": httpMessage.ERROR_MISSING_INVALID_TOKEN,
			})
			c.Abort()
		}
		exp := payload["exp"]
		id := payload["sub"]
		if int64(exp.(float64)) <= time.Now().UnixMilli() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": httpMessage.ERROR_EXPIRED_TOKEN,
			})
			c.Abort()
		}
		objectId, err := primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": httpMessage.ERROR_MISSING_INVALID_TOKEN,
			})
			c.Abort()
		}
		filterStage := bson.D{
			{
				Key:   "$match",
				Value: bson.D{{Key: "_id", Value: objectId}},
			},
		}
		var stages []bson.D
		stages = append(stages, filterStage)
		account := accountStorage.GetAccountBy(db, stages)
		if account[0].Email == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": httpMessage.ERROR_UNAUTHORIZED,
			})
			c.Abort()
		}
		c.Set("account", account[0])
		c.Next()
	}
}
