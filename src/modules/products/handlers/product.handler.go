package productHandler

import (
	"net/http"

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
		}
	}
}
