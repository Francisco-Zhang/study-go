package main

import (
	"OldPackageTest/new_helloworld/handler"
	"OldPackageTest/new_helloworld/server_proxy"
	"net"
	"net/rpc"
)

func main() {
	//1.实例化一个server
	listener, _ := net.Listen("tcp", ":1234")
	//2.注册处理逻辑handler
	server_proxy.RegisterHelloService(&handler.HelloService{}) //相当于注册HelloService.Hello
	//3.启动服务
	for {
		conn, _ := listener.Accept() //当一个新的连接进来以后，就有了一个socket的套接字
		go rpc.ServeConn(conn)
	}

	//一连串的代码大部分都是net包，好像和rpc没有关系，那么rpc可以去掉吗？
	//答案是不行，rpc 调用有几个问题需要解决 1.call id, 2.序列化和反序列化
}
