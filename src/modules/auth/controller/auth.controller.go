package auth

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Route(v1 *gin.RouterGroup, db *mongo.Client) {
	authRoute := v1.Group("/auth")
	authRoute.POST("/register", Register(db))
}
