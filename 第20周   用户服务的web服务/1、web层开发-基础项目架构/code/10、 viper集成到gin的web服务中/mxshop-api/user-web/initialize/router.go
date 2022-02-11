package initialize

import (
	"github.com/gin-gonic/gin"
	router2 "mxshop-api/user-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	//生成全局所有分组的顶层分组
	ApiGroup := Router.Group("/u/v1")
	router2.InitUserRouter(ApiGroup)

	return Router
}
