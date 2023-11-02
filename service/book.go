package service

import (
	"bookmanager-server/global"
	"bookmanager-server/model"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func GetBooks(currentPage int, pageSize int, keyWord string) (gin.H, error) {
	var books []model.Book
	var count int64
	opts := options.Find().SetSkip(int64((currentPage - 1) * pageSize)).SetLimit(int64(pageSize))
	filter := bson.M{}

	if keyWord != "" {
		regexString := fmt.Sprintf(".*%v.*", keyWord)
		filter = bson.M{
			"$or": []bson.M{
				{
					"name": primitive.Regex{Pattern: regexString},
				},
				{
					"press": primitive.Regex{Pattern: regexString},
				},
				{
					"author": primitive.Regex{Pattern: regexString},
				},
			},
		}
	}

	cursor, err := global.BookColl.Find(context.TODO(), filter, opts)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err = cursor.All(context.TODO(), &books); err != nil {
		fmt.Println(err)
		return nil, err
	}

	count, err = global.BookColl.CountDocuments(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return gin.H{
		"books": books,
		"count": count,
	}, nil
}

func AddBooks(books []model.Book) (gin.H, error) {
	var existBooks []model.Book

	// 逐个查重并添加
	for i := 0; i < len(books); i++ {
		err := global.BookColl.FindOne(context.TODO(), bson.M{
			"name":   books[i].Name,
			"author": books[i].Author,
			"press":  books[i].Press,
			"cover":  books[i].Cover,
		}).Decode(&model.Book{})

		// 找到说明存在重复的书
		if err == nil {
			existBooks = append(existBooks, books[i])
		} else {
			_, err = global.BookColl.InsertOne(context.TODO(), books[i])
			if err != nil {
				return nil, err
			}
		}
	}

	if len(existBooks) != 0 {
		return gin.H{
			"existBook": existBooks,
		}, errors.New("存在相同的书")
	} else {
		return gin.H{}, nil
	}
}

func UpdateBook(book *model.Book, changeLendStatus bool) (gin.H, error) {
	var DBBook model.Book
	if err := global.BookColl.FindOne(context.TODO(), bson.M{"_id": book.Id}).Decode(&DBBook); err != nil {
		fmt.Println(err)
		return nil, err
	}
	if book.Name != "" {
		DBBook.Name = book.Name
	}
	if book.Press != "" {
		DBBook.Press = book.Press
	}
	if len(book.Author) != 0 {
		DBBook.Author = book.Author
	}
	if book.Cover != "" {
		DBBook.Cover = book.Cover
	}
	if changeLendStatus {
		DBBook.Lend = !DBBook.Lend
		DBBook.LogDate = time.Now().Format("2006-01-02 15:04:05")
	}
	filter := bson.M{"_id": book.Id}
	update := bson.M{"$set": DBBook}
	opts := options.FindOneAndUpdate().SetReturnDocument(1)
	if err := global.BookColl.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&DBBook); err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		return gin.H{
			"book": DBBook,
		}, nil
	}
}

func DeleteBook(bookId primitive.ObjectID) (gin.H, error) {
	if _, err := global.BookColl.DeleteOne(context.TODO(), bson.M{"_id": bookId}); err != nil {
		return gin.H{
			"status": false,
		}, err
	} else {
		return gin.H{
			"status": true,
		}, nil
	}
}

func GetBook(bookId primitive.ObjectID) (gin.H, error) {
	var DBBook model.Book
	if err := global.BookColl.FindOne(context.TODO(), bson.M{"_id": bookId}).Decode(&DBBook); err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		return gin.H{
			"book": DBBook,
		}, nil
	}
}
