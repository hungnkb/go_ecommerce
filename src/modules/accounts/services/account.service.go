package accountService

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/gin-gonic/gin"
	accountModel "github.com/hungnkb/go_ecommerce/src/modules/accounts/models"
	accountStorage "github.com/hungnkb/go_ecommerce/src/modules/storages"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetList(db *mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		search := c.Query("search")
		filter := bson.M{}
		res := accountStorage.GetAccountList(db, filter, int64(page), int64(limit), search)
		if res.Error == "" {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"data":   res.Data,
			})
			return
		} else {
			c.JSON(res.HttpStatusCode, gin.H{
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
			accountStorage.InsertAccount(db, accountModel.Account{Username: petname.Name(), Password: "123123", Phone: phone, Gender: gender, Dob: primitive.DateTime(time.Now().Unix()), Email: email, Name: name})
		}
	}
	return fn
}

func CreatePermission(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []accountModel.Permission
		if error := c.ShouldBindJSON(&body); error != nil {
			c.JSON(400, gin.H{
				"message": error.Error(),
			})
			return
		}
		res := accountStorage.InsertPermissionBulk(db, body)
		if res.Error != "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  res.HttpStatusCode,
				"message": res.Error,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   res.Data,
		})
	}
}
