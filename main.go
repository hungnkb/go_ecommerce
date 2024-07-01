package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hungnkb/go_ecommerce/src/modules/accounts/accountController"
	authController "github.com/hungnkb/go_ecommerce/src/modules/auth/authControllers"
	"github.com/hungnkb/go_ecommerce/src/modules/storage"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	r := gin.Default()

	db := storage.NewMongoStorage()

	v1 := r.Group("/v1")
	accountController.Account(db, v1)
	authController.Auth(db, v1)

	PORT := os.Getenv("PORT")
	r.Run(":" + PORT)
}
