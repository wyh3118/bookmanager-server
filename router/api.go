package router

import (
	"bookmanager-server/controller"
	"github.com/gin-gonic/gin"
)

func getApiRouter(engin *gin.Engine) {
	engin.GET("/apis", controller.GetApis)     // 分页获取api列表
	engin.POST("/apis", controller.PostApis)   // 添加api
	engin.PUT("/api", controller.UpdateApi)    // 修改api
	engin.DELETE("/api", controller.DeleteApi) // 删除api
}
