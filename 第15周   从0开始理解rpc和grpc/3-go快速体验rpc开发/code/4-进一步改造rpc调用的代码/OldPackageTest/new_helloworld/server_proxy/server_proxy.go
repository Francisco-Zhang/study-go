package server_proxy

import (
	"OldPackageTest/new_helloworld/handler"
	"net/rpc"
)

type HelloService interface {
	Hello(request string, reply *string) error
}

//服务端注册如何做到解耦，我们关心的是结构体内的函数名而不是结构体类型，所以使用鸭子类型，将方法抽象成服务接口。

func RegisterHelloService(srv HelloService) error {
	return rpc.RegisterName(handler.HelloServiceName, srv)
}
