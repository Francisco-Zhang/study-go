package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	ServiceName string `mapstructure:"name"`
	Port        int    `mapstructure:"port"`
}

func main() {
	v := viper.New()
	//文件的路径如何设置 goland 运行配置里的 working directory就是工作目录
	//命令行运行，需要切换到 项目根目录，然后运行 go run viper_test/ch01/main.go
	v.SetConfigFile("viper_test/ch01/config.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	serverConfig := ServerConfig{}
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	fmt.Println(serverConfig)
	fmt.Printf("%V", v.Get("name"))
}
