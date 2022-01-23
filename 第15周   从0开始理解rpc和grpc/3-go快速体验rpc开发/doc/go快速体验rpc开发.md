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

```python
import json
import socket# 客户端 发送一个数据，再接收一个数据

request = {
    "id":0,
    "params":["imooc"],
    "method": "HelloService.Hello"
}

client = socket.create_connection(("localhost", 1234))
client.sendall(json.dumps(request).encode())

# This must loop if resp is bigger than 4K
response = client.recv(4096)
response = json.loads(response.decode())
print(response)

client.close() #关闭这个链接
```



## 3、替换rpc的传输协议为http

### server端

```go
package main

import (
	"io"
	"net/http"
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
	_ = rpc.RegisterName("HelloService", &HelloService{}) //相当于注册HelloService.Hello
	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer:     w,
		}
		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})
	http.ListenAndServe(":1234", nil)
}
```

### 客户端 Python 

安装 requests 包：pip install requests

```python
import requests

request = {
    "id":0,
    "params":["imooc"],
    "method": "HelloService.Hello"
}

rsp= requests.post("http://localhost:1234/jsonrpc",json=request)
print(rsp.text)
```



## 4、进一步改造rpc调用的代码

1、服务名称由原来的硬编码改为使用变量

2、只想写业务逻辑，不想关注每个函数的名称，所以需要自定义一个client 的代理，实现对rpc方法的封装。

3、服务端的业务逻辑抽离，放到handler中。服务端的服务注册放入到server_proxy中。这样新增加业务时，就不需要修改server.go 和 client.go 了。

4、服务端注册如何做到解耦，我们关心的是结构体内的函数名而不是结构体类型，所以使用鸭子类型，将方法抽象成服务接口。传入的struct只要包含Hello方法就行，所以将接收参数抽象成interface。

以上这些概念在grpc中都有对应。发自灵魂的拷问，处理业务的 server_proxy 和 client_proxy 能不能自动生成？为多种语言生成？

这两个要求都能满足，这个就是 protobuf + grpc 。

