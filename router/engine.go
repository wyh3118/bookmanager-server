package router

import (
	"bookmanager-server/middleware"
	"github.com/gin-gonic/gin"
)

func GetEngine() *gin.Engine {
	engine := gin.Default()

	// 配置跨域,使用插件的默认配置
	engine.Use(middleware.Cors())
	getUserRouter(engine)
	// 使用用户认证中间件
	//engine.Use(middleware.AuthMiddleware(), middleware.ApiAuth())
	engine.Use(middleware.AuthMiddleware())
	getApiRouter(engine)
	getRoleRouter(engine)
	getBookRouter(engine)
	getAdminOperateUserRouter(engine)
	return engine
}
