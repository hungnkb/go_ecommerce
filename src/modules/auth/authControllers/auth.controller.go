package authController

import (
	"github.com/gin-gonic/gin"
	authServices "github.com/hungnkb/go_ecommerce/src/modules/auth/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func Auth(db *mongo.Client, r *gin.RouterGroup) {
	authRoute := r.Group("/auth")
	authRoute.POST("/register", authServices.Register(db))
	authRoute.POST("/login", authServices.LoginByPassword(db))
}
