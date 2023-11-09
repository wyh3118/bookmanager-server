package initialize

import (
	"bookmanager-server/global"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func MongoInit() {
	if global.MongoClient == nil {
		global.MongoClient = getMongoClient("mongodb://admin:123456@43.138.56.2")
	}

	BookManager := global.MongoClient.Database("BookManager")
	{
		global.UserColl = BookManager.Collection("user")
		global.BookColl = BookManager.Collection("book")
		global.ApiColl = BookManager.Collection("api")
		global.RoleColl = BookManager.Collection("role")
	}
}

func getMongoClient(url string) *mongo.Client {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connect mongodb!")
	return client
}
