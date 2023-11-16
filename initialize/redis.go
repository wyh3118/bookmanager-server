package initialize

import (
	"bookmanager-server/global"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

func RedisInit() {
	global.Redis = redis.NewClient(&redis.Options{
		Addr:     "43.138.73.97:6379", // 数据库地址
		Username: "default",
		Password: "redis654321",
		DB:       0,
		PoolSize: 20,
	})

	_, err := global.Redis.Ping(context.TODO()).Result()
	if err == nil {
		fmt.Println("redis connect!")
	} else {
		log.Fatalln(err)
	}
}
