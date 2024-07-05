package productController

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Product(db *mongo.Client, r *gin.RouterGroup) {
	productRoute := r.Group("/products")
	productRoute.POST("/")
	productRoute.GET("/")
	productRoute.GET("/:id")
	productRoute.PUT("/:id")
	productRoute.DELETE("/")
}
