package middleware

import (
	"bookmanager-server/global"
	"bookmanager-server/model"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func ApiAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		method, url := c.Request.Method, c.Request.URL.Path
		// 查找是否有此api
		err := global.ApiColl.FindOne(context.TODO(), bson.M{
			"url":    url,
			"method": method,
		}).Decode(&model.Api{})
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "没有此api",
			})
			c.Abort()
			return
		}

		// 根据用户ID获取用户的角色信息,再根据角色信息判断是否有权访问
		var user model.User
		_ = global.UserColl.FindOne(context.TODO(), bson.M{"_id": c.MustGet("Id").(primitive.ObjectID)}).Decode(&user)
		flag := false
	label:
		for _, roleId := range user.Roles {
			var role model.Role
			_ = global.RoleColl.FindOne(context.TODO(), bson.M{"_id": roleId}).Decode(&role)
			fmt.Println(role)
			for _, apiId := range role.Apis {
				var api model.Api
				_ = global.ApiColl.FindOne(context.TODO(), bson.M{"_id": apiId}).Decode(&api)
				if method == api.Method && url == api.Url {
					flag = true
					break label
				}
			}
		}

		if !flag {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "无权访问",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
