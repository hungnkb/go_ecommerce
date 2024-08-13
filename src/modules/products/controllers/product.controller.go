package productController

import (
	"github.com/gin-gonic/gin"
	middleware "github.com/hungnkb/go_ecommerce/src/middlewares"
	productHandler "github.com/hungnkb/go_ecommerce/src/modules/products/handlers"
	"go.mongodb.org/mongo-driver/mongo"
)

func Product(db *mongo.Client, r *gin.RouterGroup) {
	productRoute := r.Group("/products")
	productRoute.POST("/", middleware.AuthGuard(db), productHandler.Create(db))
	productRoute.GET("/", productHandler.GetList(db))
	productRoute.GET("/:slug", productHandler.GetProductBySlug(db))
	productRoute.PUT("/:id")
	productRoute.DELETE("/")

	productRoute.POST("/attributes", middleware.AuthGuard(db), productHandler.CreateAttributeBulk(db))
	productRoute.GET("/attributes", middleware.AuthGuard(db), productHandler.CreateAttributeBulk(db))

	productRoute.GET("/categories", productHandler.GetListCategory(db))
	productRoute.POST("/categories", middleware.AuthGuard(db), productHandler.CreateCategory(db))
	productRoute.GET("/categories/:id")
}
