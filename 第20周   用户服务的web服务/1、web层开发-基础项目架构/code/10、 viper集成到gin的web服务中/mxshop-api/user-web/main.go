package main

import (
	"fmt"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
)

func main() {
	//1. 初始化logger
	initialize.InitLogger()
	//2. 初始化配置文件
	initialize.InitConfig()

	//3. 初始化routers
	Router := initialize.Routers()
	//使用全局的log,缩短代码
	/*
		1. S()可以获取一个全局的sugar，这个sugar什么都没做，需要我们自己设置一个全局的logger
			logger, _ := zap.NewDevelopment()
			zap.ReplaceGlobals(logger)
		2. 日志是分级别的，debug， info ， warn， error， fetal,生成环境不打印debug
		3. S函数和L函数很有用， 提供了一个全局的安全访问logger的途径
	*/
	zap.S().Infof("启动服务器,端口：%d", global.ServerConfig.Port)

	err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port))
	if err != nil {
		zap.S().Panic("启动失败", err.Error())
	}
}
