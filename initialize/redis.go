package initialize

import (
	"bookmanager-server/global"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
)

func RedisInit() {
	global.Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")), // 数据库地址
		Username: viper.GetString("redis.username"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.poolsize"),
	})

	_, err := global.Redis.Ping(context.TODO()).Result()
	if err == nil {
		fmt.Println("redis connect!")
	} else {
		log.Fatalln(err)
	}
}
