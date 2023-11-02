package controller

import (
	"bookmanager-server/model"
	"bookmanager-server/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

func GetBooks(c *gin.Context) {
	var data struct {
		CurrentPage int    `form:"currentPage" binding:"required"`
		PageSize    int    `form:"pageSize" binding:"required"`
		KeyWord     string `form:"keyWord"`
	}

	if err := c.ShouldBind(&data); err != nil {
		fmt.Println("data:", data)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}

	if result, err := service.GetBooks(data.CurrentPage, data.PageSize, data.KeyWord); err != nil {
		c.JSON(http.StatusBadRequest, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func AddBooks(c *gin.Context) {
	var data struct {
		BookList []model.Book `json:"bookList" binding:"required"`
	}
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}
	books := data.BookList
	for i := 0; i < len(books); i++ {
		books[i].Lend = false
		books[i].LogDate = time.Now().Format("2006-01-02 15:04:05")
	}

	if result, err := service.AddBooks(books); err != nil {
		c.JSON(http.StatusBadRequest, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func UpdateBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBind(&book); err != nil {
		fmt.Println(book)
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
	}

	if bookId, err := primitive.ObjectIDFromHex(c.Query("Id")); err != nil {
		fmt.Println(book)
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
	} else {
		book.Id = bookId
	}

	var changeLendStatus bool
	if c.Query("changeLendStatus") == "" {
		changeLendStatus = false
	} else {
		changeLendStatus = true
	}

	if result, err := service.UpdateBook(&book, changeLendStatus); err != nil {
		c.JSON(http.StatusBadRequest, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func DeleteBook(c *gin.Context) {
	bookID, err := primitive.ObjectIDFromHex(c.Query("Id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
	}

	if result, err := service.DeleteBook(bookID); err != nil {
		c.JSON(http.StatusBadRequest, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func GetBook(c *gin.Context) {
	bookID, err := primitive.ObjectIDFromHex(c.Query("Id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
	}

	if result, err := service.GetBook(bookID); err != nil {
		c.JSON(http.StatusBadRequest, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}
