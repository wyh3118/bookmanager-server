package router

import (
	"bookmanager-server/controller"
	"github.com/gin-gonic/gin"
)

func getBookRouter(engine *gin.Engine) {
	engine.GET("/books", controller.GetBooks)
	engine.POST("/books", controller.AddBooks)
	engine.GET("/book", controller.GetBook)
	engine.PUT("/book", controller.UpdateBook)
	engine.DELETE("/book", controller.DeleteBook)
}
