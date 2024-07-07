package accountController

import (
	"github.com/gin-gonic/gin"
	middleware "github.com/hungnkb/go_ecommerce/src/middlewares"
	accountService "github.com/hungnkb/go_ecommerce/src/modules/accounts/accountHandler"
	"go.mongodb.org/mongo-driver/mongo"
)

func Account(db *mongo.Client, r *gin.RouterGroup) {
	accountRoute := r.Group("/accounts")
	accountRoute.GET("/", middleware.AuthGuard(db), accountService.GetList(db))
	accountRoute.POST("/mock", middleware.AuthGuard(db), accountService.MockAccount(db))
	accountRoute.POST("/permissions", accountService.CreatePermission(db))
}
