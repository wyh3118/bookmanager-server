package router

import (
	"bookmanager-server/controller"
	"github.com/gin-gonic/gin"
)

func getRoleRouter(engin *gin.Engine) {
	engin.GET("/roles", controller.GetRoles)
	engin.POST("roles", controller.PostRoles)
	engin.PUT("/role", controller.UpdateRole)
	engin.DELETE("/role", controller.DeleteRole)
}
