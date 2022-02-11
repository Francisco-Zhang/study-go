package main

import (
	"fmt"
	"go.uber.org/zap"
	"mxshop-api/user-web/initialize"
)

func main() {
	//初始化routers
	Router := initialize.Routers()
	zap.NewProduction()

	port := 8021
	err := Router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}
}
