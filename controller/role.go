package controller

import (
	"bookmanager-server/model"
	"bookmanager-server/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func GetRoles(c *gin.Context) {
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

	if result, err := service.GetRoles(data.CurrentPage, data.PageSize, data.KeyWord); err != nil {
		c.JSON(http.StatusBadRequest, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func PostRoles(c *gin.Context) {
	var data struct {
		Roles []model.Role `json:"roles" bson:"roles"`
	}

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}

	if result, err := service.PostRoles(data.Roles); err != nil {
		c.JSON(http.StatusBadRequest, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func UpdateRole(c *gin.Context) {
	role := model.Role{}
	if err := c.Bind(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}

	Id, err := primitive.ObjectIDFromHex(c.Query("Id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	} else {
		role.Id = Id
	}

	if result, err := service.UpdateRole(&role); err != nil {
		c.JSON(http.StatusBadRequest, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func DeleteRole(c *gin.Context) {
	roleId, err := primitive.ObjectIDFromHex(c.Query("Id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
	}

	if result, err := service.DeleteRole(roleId); err != nil {
		c.JSON(http.StatusBadRequest, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}
