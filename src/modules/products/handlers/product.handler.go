package productHandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	accountModel "github.com/hungnkb/go_ecommerce/src/modules/accounts/models"
	productStorage "github.com/hungnkb/go_ecommerce/src/modules/products/storages"
	"go.mongodb.org/mongo-driver/mongo"
)

func Create(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentAccount, _ := c.Get("account")
		account := currentAccount.(accountModel.Account)
		var input productStorage.ProductInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		response := productStorage.InsertProduct(db, input, account)
		if response.Error != "" {
			c.JSON(response.HttpStatusCode, gin.H{"message": response.Error})
			return
		}
		c.JSON(response.HttpStatusCode, gin.H{"data": response.Data})
	}
}

func GetList(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.Query("page")
		limit := c.Query("limit")
		keywords := c.Query("keywords")
		if page == "" {
			page = "1"
		}
		if limit == "" {
			limit = "24"
		}
		limitInt, _ := strconv.ParseInt(limit, 10, 64)
		pageInt, _ := strconv.ParseInt(page, 10, 64)
		// if sort == "" {
		// 	sort = ""
		// }
		response := productStorage.ProductGetList(db, pageInt, limitInt, keywords)
		if response.Error != "" {
			c.JSON(response.HttpStatusCode, gin.H{"message": response.Error})
			return
		}
		c.JSON(response.HttpStatusCode, gin.H{
			"data": response.Data,
		})
	}
}
