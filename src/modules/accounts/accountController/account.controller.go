package accountController

import (
	"github.com/gin-gonic/gin"
	service "github.com/hungnkb/go_ecommerce/src/modules/accounts/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func Account(db *mongo.Client, r *gin.RouterGroup) {
	accountRoute := r.Group("/accounts")
	accountRoute.GET("", service.GetList(db))
}
