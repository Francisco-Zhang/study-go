## 1、Go 内置 rpc 快速开发

### server端

```go
package main

import (
	"net"
	"net/rpc"
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
	conn, _ := listener.Accept() //当一个新的连接进来以后，就有了一个socket的套接字
	rpc.ServeConn(conn)

	//一连串的代码大部分都是net包，好像和rpc没有关系，那么rpc可以去掉吗？
	//答案是不行，rpc 调用有几个问题需要解决 1.call id, 2.序列化和反序列化
}

```

### client端

```go
package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic("连接失败")
	}
	var reply string
	err = client.Call("HelloService.Hello", "tom", &reply)
	if err != nil {
		panic("调用失败")
	}
	fmt.Println(reply)
}
```



## 2、替换rpc的序列化协议为json

client 的调用使用起来并不方便。理想的调用是 client.Hello()。是否可以跨语言调用呢，需要考虑两点

1. go 语言的 rpc 的序列化和反序列化协议是什么 （Gob 协议）
2. 能否替换成常用的序列化

### server端

```go
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
```

### client端

```go
package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		panic("连接失败")
	}
	var reply string
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	err = client.Call("HelloService.Hello", "tom", &reply)
	if err != nil {
		panic("调用失败")
	}
	fmt.Println(reply)
}

```

### 使用 Python 进行客户端改造

首先要知道 client 发啥数据的 json 格式

```json
{"method":"HelloService.Hello","params":["hello"],"id":0}
```

