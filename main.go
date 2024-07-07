package main

import (
	"github.com/gin-gonic/gin"
	Config "github.com/hungnkb/go_ecommerce/src/config"
	accountController "github.com/hungnkb/go_ecommerce/src/modules/accounts/controllers"
	authController "github.com/hungnkb/go_ecommerce/src/modules/auth/controllers"
	productController "github.com/hungnkb/go_ecommerce/src/modules/products/controllers"
	accountStorage "github.com/hungnkb/go_ecommerce/src/modules/storage"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	gin.ForceConsoleColor()
	r := gin.Default()

	db := accountStorage.NewMongoStorage()
	api := r.Group("/api")
	v1 := api.Group("/v1")

	accountController.Account(db, v1)
	authController.Auth(db, v1)
	productController.Product(db, v1)

	PORT := Config.Get().Port
	r.Run(":" + PORT)
}
