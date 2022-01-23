package client_proxy

import (
	"OldPackageTest/new_helloworld/handler"
	"net/rpc"
)

// HelloServiceStub 取名Stub 是为了和 grpc 命名保持一致
type HelloServiceStub struct {
	*rpc.Client
}

//go 语言中没有类、对象，就意味着没有初始化方法
//所以定义一个New开头的初始化函数，在Go语言中非常常见。

func NewHelloServiceClient(protol, address string) HelloServiceStub {
	client, err := rpc.Dial(protol, address)
	if err != nil {
		panic("connect failed")
	}
	return HelloServiceStub{client}
}

func (c *HelloServiceStub) Hello(request string, reply *string) error {
	err := c.Call(handler.HelloServiceName+".Hello", request, reply)
	if err != nil {
		return err
	}
	return nil
}
