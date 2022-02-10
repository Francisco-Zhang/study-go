package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func MyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Set("example", "123456")
		//让原本改执行的逻辑继续执行
		c.Next()

		end := time.Since(t)
		fmt.Printf("耗时:%V\n", end)
		status := c.Writer.Status()
		fmt.Println("状态", status)
	}
}

func Hook404() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
		status := c.Writer.Status()
		if status == 404 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "页面找不到",
			})
		}
	}
}

func TokenRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		for k, v := range c.Request.Header {
			if k == "X-Token" {
				token = v[0]
			}
		}
		if token != "bobby" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "未登录",
			})
			// 挑战，我们会有一个疑问，为什么 return 都阻止不了后续逻辑的执行。
			c.Abort() //此处使用 return 仍然能正常返回接口响应。响应结果有两个json，一个是未登录，一个是 pong
		}
		c.Next()
	}
}

func main() {
	router := gin.Default()
	//使用logger和recovery中间件 全局所有
	router.Use(Hook404())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Run(":8083")
}
