package authService

import (
	"github.com/gin-gonic/gin"
	"github.com/hungnkb/go_ecommerce/src/modules/accounts/models"
	"github.com/hungnkb/go_ecommerce/src/modules/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(db *mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var body models.Account
		if error := c.ShouldBind(&body); error != nil {
			c.JSON(400, gin.H{
				"message": error.Error(),
			})
		}

		data := storage.InsertAccount(db, body)
		if data.Data != nil {
			c.JSON(200, gin.H{
				"status": 200,
				"data":   data,
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
