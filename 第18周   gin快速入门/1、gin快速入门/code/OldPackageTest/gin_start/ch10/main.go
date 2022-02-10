package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//优雅退出 就是当我们关闭程序的时候，保证程序应该做完一些后续处理再退出，而不是突然终止造成数据丢失。
	//微服务 启动之前或启动之后都会做一件事：将当前服务的ip地址和端口号注册到注册中心
	//没有优雅退出，导致程序停止了之后没有告知注册中心。
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	go func() {
		router.Run(":8083")
	}()

	//如果想要接收到信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //ctl+c 和 kill 信号

	<-quit
	//处理后续逻辑，完成优雅退出
	fmt.Println("关闭服务中.....")
	fmt.Println("注销服务.....")
}
