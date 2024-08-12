package productHandler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	productModel "github.com/hungnkb/go_ecommerce/src/modules/products/models"
	productStorage "github.com/hungnkb/go_ecommerce/src/modules/products/storages"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateCategory(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input productModel.Category
		c.ShouldBindJSON(&input)
		response := productStorage.InsertCategory(db, input)
		if response.Error != "" {
			c.JSON(response.HttpStatusCode, gin.H{"message": response.Error})
			return
		}
		c.JSON(response.HttpStatusCode, gin.H{"data": response.Data})
	}
}

func GetListCategory(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		parentId := c.Query("parentId")
		page := c.Query("page")
		limit := c.Query("limit")
		if page == "" {
			page = "1"
		}
		if limit == "" {
			limit = "1000"
		}
		pageInt, _ := strconv.Atoi(page)
		limitInt, _ := strconv.Atoi(limit)
		response := productStorage.GetListCategory(db, pageInt, limitInt, parentId)
		if response.Error != "" {
			c.JSON(response.HttpStatusCode, gin.H{"message": response.Error})
			return
		}
		c.JSON(response.HttpStatusCode, gin.H{"data": response.Data})
	}
}
