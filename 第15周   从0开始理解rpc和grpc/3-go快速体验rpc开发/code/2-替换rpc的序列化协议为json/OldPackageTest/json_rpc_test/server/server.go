package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloService struct {
}

func (s *HelloService) Hello(request string, reply *string) error {
	*reply = "hello," + request
	return nil
}

func main() {
	//1.实例化一个server
	listener, _ := net.Listen("tcp", ":1234")
	//2.注册处理逻辑handler
	_ = rpc.RegisterName("HelloService", &HelloService{}) //相当于注册HelloService.Hello
	//3.启动服务
	for {
		conn, _ := listener.Accept()                    //当一个新的连接进来以后，就有了一个socket的套接字
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn)) //使用协程，同时处理多个连接
	}

}
