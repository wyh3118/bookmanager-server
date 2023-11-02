package controller

import (
	"bookmanager-server/model"
	"bookmanager-server/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

func checkMethod(method string) bool {
	methods := []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
		"PATCH",
	}
	for i := 0; i < len(methods); i++ {
		if method == methods[i] {
			return true
		}
	}
	return false
}

func GetApis(c *gin.Context) {
	var data struct {
		CurrentPage int    `form:"currentPage" binding:"required"`
		PageSize    int    `form:"pageSize" binding:"required"`
		KeyWord     string `form:"keyWord"`
	}

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}

	if result, err := service.GetApis(data.CurrentPage, data.PageSize, data.KeyWord); err != nil {
		c.JSON(http.StatusBadRequest, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func PostApis(c *gin.Context) {
	data := struct {
		Apis []model.Api `json:"apis" bson:"apis"`
	}{}

	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}

	// 将method统一转换为大写并验证
	for i := 0; i < len(data.Apis); i++ {
		data.Apis[i].Method = strings.ToUpper(data.Apis[i].Method)
		if !checkMethod(data.Apis[i].Method) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "method参数错误",
			})
			return
		}
	}

	if result, err := service.PostApis(data.Apis); err != nil {
		c.JSON(http.StatusBadRequest, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func UpdateApi(c *gin.Context) {
	var api model.Api
	if err := c.ShouldBind(&api); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}

	if apiId, err := primitive.ObjectIDFromHex(c.Query("Id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	} else {
		api.Id = apiId
	}

	api.Method = strings.ToUpper(api.Method)
	if !checkMethod(api.Method) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "method参数错误",
		})
		return
	}

	if result, err := service.UpdateApi(&api); err != nil {
		c.JSON(http.StatusBadRequest, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func DeleteApi(c *gin.Context) {
	apiId, err := primitive.ObjectIDFromHex(c.Query("Id"))
	fmt.Println(apiId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}

	if result, err := service.DeleteApi(apiId); err != nil {
		c.JSON(http.StatusBadRequest, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}
