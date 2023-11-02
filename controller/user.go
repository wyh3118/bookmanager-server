package controller

import (
	"bookmanager-server/model"
	"bookmanager-server/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"regexp"
)

func Login(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
	}

	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}

	if result, err := service.GetUser(&user); err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, result)
	}
}

func Register(c *gin.Context) {
	user := &model.User{}
	if err := c.Bind(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}

	// 校验邮箱格式
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	ok := reg.MatchString(user.Email)
	if ok == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "邮箱格式错误",
		})
		return
	}

	if result, err := service.CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, result)
	} else {
		c.JSON(http.StatusCreated, result)
	}
}

func UpdateUserInfo(c *gin.Context) {
	user := &model.User{}
	if err := c.Bind(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}

	user.Id = c.MustGet("Id").(primitive.ObjectID)

	if result, err := service.UpdateUser(user); err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, result)
	}
}

func UpdateUserByAdmin(c *gin.Context) {
	var user model.User
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
	}

	if result, err := service.UpdateUserByAdmin(&user); err != nil {
		c.JSON(http.StatusBadRequest, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}
