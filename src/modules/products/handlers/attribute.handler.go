package productHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	accountModel "github.com/hungnkb/go_ecommerce/src/modules/accounts/models"
	productModel "github.com/hungnkb/go_ecommerce/src/modules/products/models"
	productStorage "github.com/hungnkb/go_ecommerce/src/modules/products/storages"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateAttributeBulk(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input []productModel.ProductAttribute
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		currentAccount, _ := c.Get("account")
		res := productStorage.InsertAttributeBulk(db, input, currentAccount.(accountModel.Account))
		if res.Error != "" {
			c.JSON(res.HttpStatusCode, gin.H{"message": res.Error})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   res.Data,
		})
	}

}
