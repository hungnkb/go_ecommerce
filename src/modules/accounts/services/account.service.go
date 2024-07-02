package service

import (
	"math/rand"
	"strconv"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/gin-gonic/gin"
	"github.com/hungnkb/go_ecommerce/src/modules/accounts/models"
	"github.com/hungnkb/go_ecommerce/src/modules/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetList(db *mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		filter := bson.D{{}}
		res := storage.GetAccountList(db, filter, int64(page), int64(limit))
		if res.Error == "" {
			c.JSON(200, gin.H{
				"message": "ping",
				"data":    res.Data,
			})
			return
		} else {
			c.JSON(200, gin.H{
				"message": res.Error,
				"data":    nil,
			})
			return
		}
	}
	return fn
}

func MockAccount(db *mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.Query("limit"))
		if limit == 0 {
			limit = 10
		}
		for range limit {
			phone := "0"
			for v := range 10 {
				if v > 0 {
					phone += strconv.Itoa(rand.Intn(10))
				}
			}
			gender := "male"
			name := petname.Name()
			var email string
			isEven := time.Now().Local().UnixMicro()%2 == 0
			if isEven {
				gender = "female"
				email = name + "@yahoo.com"
			} else {
				email = name + "@gmail.com"
			}
			storage.InsertAccount(db, models.Account{Username: petname.Name(), Password: "123123", Phone: phone, Gender: gender, Dob: primitive.DateTime(time.Now().Unix()), Email: email, Name: name})
		}
	}
	return fn
}
