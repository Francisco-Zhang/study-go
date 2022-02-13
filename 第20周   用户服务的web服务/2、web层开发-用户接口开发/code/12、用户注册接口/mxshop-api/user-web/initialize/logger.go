package initialize

import "go.uber.org/zap"

func InitLogger() {
	//生成环境不打印debug, NewProduction 日志格式是 json 格式的
	//logger, _ := zap.NewProduction()
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
