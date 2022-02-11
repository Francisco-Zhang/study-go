## 1、 新建项目和目录结构构建

1. 新建项目mxshop-api
2. 新建目录user-web
3. user-web下面新建子目录 api、config、forms、global、initialize、middlewares、proto、router、utils、validator等目录
4. 项目根目录新建启动文件 main.go

## 2、 go高性能日志库 - zap使用

zap性能非常高，很多公司普遍使用。

```shell
go get -u go.uber.org/zap
```

```go
package main

import (
	"go.uber.org/zap"
)

func main()  {
	logger, _ := zap.NewProduction() //生产环境
	//logger, _ := zap.NewDevelopment()//开发环境 、log level 等配置不同
	defer logger.Sync() // flushes buffer, if any
	url := "https://imooc.com"
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
}
```

Zap提供了两种类型的日志记录器—`Sugared Logger`和`Logger`。

在性能很好但不是很关键的上下文中，使用`SugaredLogger`。它比其他结构化日志记录包快4-10倍，并且支持结构化和printf风格的日志记录。

在每一微秒和每一次内存分配都很重要的上下文中，使用`Logger`。它甚至比`SugaredLogger`更快，内存分配次数也更少，但它只支持强类型的结构化日志记录

```go
//第二种日志类型使用。指定了类型，减少了反射，性能高，但是不好用。
logger.Info("failed to fetch URL",
			zap.String("url", url),
            zap.Int("numbers", 4))
```

## 3、 zap的文件输出

```go
package main

import (
	"go.uber.org/zap"
	"time"
)

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./myproject.log", "stdout", //输出到文件和控制台
	}
	return cfg.Build()
}

func main() {
	//logger, _ := zap.NewProduction()
	logger, err := NewLogger()
	if err != nil {
		panic(err)
	}
	su := logger.Sugar()
	defer su.Sync()
	url := "https://imooc.com"
	su.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}
```

## 4 、集成路由初始到gin的启动过程 

## 5 、集成zap初始到gin的启动过程

1. S()可以获取一个全局的sugar，这个sugar什么都没做，需要我们自己设置一个全局的logger
   logger, _ := zap.NewDevelopment()
   zap.ReplaceGlobals(logger)
2. 日志是分级别的，debug， info ， warn， error， fetal
3. S函数和L函数很有用， 提供了一个全局的安全访问logger的途径，S()内部有锁。

```go
//生成环境不打印debug, NewProduction 日志格式是 json 格式的
//logger, _ := zap.NewProduction()
logger, _ := zap.NewDevelopment()
zap.ReplaceGlobals(logger)

port := 8021
zap.S().Infof("启动服务器,端口：%d", port) //这种写法可以缩短代码长度。
```

## 6 、gin调用grpc服务-1

## 7 、gin调用grpc服务-2