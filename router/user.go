package router

import (
	"bookmanager-server/controller"
	"bookmanager-server/middleware"
	"github.com/gin-gonic/gin"
)

func getUserRouter(engine *gin.Engine) {
	userGroup := engine.Group("/user")
	{
		userGroup.POST("/login", controller.Login)
		userGroup.POST("/register", controller.Register)
		userGroup.PUT("", middleware.AuthMiddleware(), controller.UpdateUserInfo)
	}
}

func getAdminOperateUserRouter(engine *gin.Engine) {
	engine.PUT("/user/updateByAdmin", controller.UpdateUserByAdmin)
}
