package initialize

import (
	"bookmanager-server/global"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func MongoInit() {
	if global.MongoClient == nil {
		host := viper.GetString("mongo.host")
		port := viper.GetString("mongo.port")
		username := viper.GetString("mongo.username")
		password := viper.GetString("mongo.password")
		url := fmt.Sprintf("mongodb://%s:%s@%s:%s/", username, password, host, port)
		global.MongoClient = getMongoClient(url)
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
	clientOptions := options.Client().ApplyURI(url)

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
