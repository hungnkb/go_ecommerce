package accountController

import (
	"github.com/gin-gonic/gin"
	middleware "github.com/hungnkb/go_ecommerce/src/middlewares"
	service "github.com/hungnkb/go_ecommerce/src/modules/accounts/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func Account(db *mongo.Client, r *gin.RouterGroup) {
	accountRoute := r.Group("/accounts")
	accountRoute.GET("/", middleware.AuthGuard(db), service.GetList(db))
	accountRoute.POST("/mock", middleware.AuthGuard(db), service.MockAccount(db))
}
