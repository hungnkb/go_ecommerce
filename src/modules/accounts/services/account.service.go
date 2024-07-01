package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hungnkb/go_ecommerce/src/modules/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetList(db *mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		data := storage.GetAccountList(db)
		if data != nil {
			c.JSON(200, gin.H{
				"message": "ping",
				"data":    data,
			})
		} else {
			c.JSON(200, gin.H{
				"message": "pong",
				"data":    nil,
			})
		}
	}
	return fn
}
