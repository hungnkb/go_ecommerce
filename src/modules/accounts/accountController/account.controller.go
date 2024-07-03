package accountController

import (
	"github.com/gin-gonic/gin"
	middleware "github.com/hungnkb/go_ecommerce/src/middlewares"
	accounts "github.com/hungnkb/go_ecommerce/src/modules/accounts/accountService"
	"go.mongodb.org/mongo-driver/mongo"
)

func Account(db *mongo.Client, r *gin.RouterGroup) {
	accountRoute := r.Group("/accounts")
	accountRoute.GET("/", middleware.AuthGuard(db), accounts.GetList(db))
	accountRoute.POST("/mock", middleware.AuthGuard(db), accounts.MockAccount(db))
	accountRoute.POST("/permissions", accounts.CreatePermission(db))
}
