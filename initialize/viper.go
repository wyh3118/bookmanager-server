package initialize

import (
	"fmt"
	"github.com/spf13/viper"
)

func ViperInit() {
	viper.SetConfigFile("./config.yaml")
	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {             // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 监控配置文件变化
	viper.WatchConfig()

	fmt.Println("viper running")
}
