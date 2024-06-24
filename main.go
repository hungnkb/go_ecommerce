package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hungnkb/go_ecommerce/src/modules/accounts/models"
	"github.com/hungnkb/go_ecommerce/src/modules/storage"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	r := gin.Default()

	db := storage.NewMongoStorage()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/accounts", func(c *gin.Context) {
		data := storage.GetAccountList(db)
		if data != nil {
			c.JSON(200, gin.H{
				"message": "ping",
				"data":    data,
			})
		}
		c.JSON(200, gin.H{
			"message": "pong",
			"data":    nil,
		})
	})
	r.POST("/accounts", func(c *gin.Context) {
		var body models.Account
		if error := c.ShouldBind(&body); error != nil {
			c.JSON(400, gin.H{
				"message": error.Error(),
			})
		}
		data := storage.InsertAccount(db, body)
		if db != nil {
			c.JSON(200, gin.H{
				"data": data,
			})
		}
		c.JSON(400, gin.H{
			"message": "error",
		})
	})

	PORT := os.Getenv("PORT")
	r.Run(":" + PORT)
}
