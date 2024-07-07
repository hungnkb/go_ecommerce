package authController

import (
	"github.com/gin-gonic/gin"
	"github.com/hungnkb/go_ecommerce/src/modules/auth/authHandler"
	"go.mongodb.org/mongo-driver/mongo"
)

func Auth(db *mongo.Client, r *gin.RouterGroup) {
	authRoute := r.Group("/auth")
	authRoute.POST("/register", authHandler.Register(db))
	authRoute.POST("/login", authHandler.LoginByPassword(db))
}
