## 1、什么是grpc和protobuf

### grpc

gRPC 是一个高性能、开源和通用的 RPC 框架，面向移动和 HTTP/2 设计。目前提供 C、Java 和 Go 语言版本，分别是：[grpc](https://github.com/grpc/grpc), [grpc-java](https://github.com/grpc/grpc-java), [grpc-go](https://github.com/grpc/grpc-go). 其中 C 版本支持 [C](https://github.com/grpc/grpc), [C++](https://github.com/grpc/grpc/tree/master/src/cpp), [Node.js](https://github.com/grpc/grpc/tree/master/src/node), [Python](https://github.com/grpc/grpc/tree/master/src/python), [Ruby](https://github.com/grpc/grpc/tree/master/src/ruby), [Objective-C](https://github.com/grpc/grpc/tree/master/src/objective-c), [PHP](https://github.com/grpc/grpc/tree/master/src/php) 和 [C#](https://github.com/grpc/grpc/tree/master/src/csharp) 支持.

![1](img/1.png)

protobuf只是一个协议，基于这个协议可以自己开发一个rpc框架。grpc 也是使用的这个protobuf协议。

### protobuf

java中的dubbo 使用了 dubbo/rmi/hessian messagepack 等协议，如果你懂了协议完全有能力自己去实现一个协议

- 习惯用 `Json、XML` 数据存储格式的你们，相信大多都没听过`Protocol Buffer（缩写为protobuf）`
- `Protocol Buffer` 其实 是 `Google`出品的一种轻量 & 高效的结构化数据存储格式，性能比 `Json、XML` 真的强！太！多！

- protobuf经历了protobuf2和protobuf3，pb3比pb2简化了很多，目前主流的版本是pb3



![2](img/2.png)

如果是服务直接内部调用使用protobuf会比较好一些，如果是作为开发接口供客户端调用使用json会比较好，因为json没有加密，容易看懂。



## 2、grpc开发环境的搭建

### 1. 下载工具

https://github.com/protocolbuffers/protobuf/releases

如果觉得下载较慢可以点击这里下载：

[📎protoc-3.13.0-win64.zip](https://www.yuque.com/attachments/yuque/0/2020/zip/159615/1603012438943-0f20e6d0-f381-4dc7-a99d-2a77031a03b1.zip)

[📎protoc-3.13.0-linux-x86_64.zip](https://www.yuque.com/attachments/yuque/0/2020/zip/159615/1603012438961-8d1df617-b453-4934-8ebe-262e6c3df02d.zip)

下载完成后解压后记得将可执行文件protoc.exe 路径添加到环境变量中



### 2. 下载go的依赖包

```shell
go get github.com/golang/protobuf/protoc-gen-go
```



protoc是protobuf文件（.proto）的编译器，可以借助这个工具把 .proto 文件转译成各种编程语言对应的源码，包含数据类型定义、调用接口等。

protoc在设计上把protobuf和不同的语言解耦了，底层用c++来实现protobuf结构的存储，然后通过插件的形式来生成不同语言的源码。可以把protoc的编译过程分成简单的两个步骤：

1）解析.proto文件，转译成protobuf的原生数据结构在内存中保存；

2）把protobuf相关的数据结构传递给相应语言的编译插件，由插件负责根据接收到的protobuf原生结构渲染输出特定语言的模板。

源码中包含的插件有 csharp、java、js、objectivec、php、python、ruby等多种。

protoc-gen-go是protobuf**编译插件**系列中的Go版本。从上一小节知道原生的protoc并不包含Go版本的插件，不过可以在github上发现专门的代码库



由于protoc-gen-go是Go写的，所以安装它变得很简单，只需要运行 `go get -u github.com/golang/protobuf/protoc-gen-go`，便可以在$GOPATH/bin目录下发现这个工具。至此，就可以通过下面的命令来使用protoc-gen-go了。

```shell
protoc --go_out=output_directory input_directory/file.proto
```

其中"--go_out="表示生成Go文件，protoc会自动寻找GOPATH中的protoc-gen-go执行文件。



### 3. proto文件

```go
syntax = "proto3";
option go_package = ".;proto";
service Greeter {
  rpc SayHello (HelloRequest) returns (HelloReply);
}

message HelloRequest {
  string name = 1; //1是编号不是值
}

message HelloReply {
  string message = 1;
}
```

### 4. 生成go文件

```shell
protoc -I . goods.proto --go_out=plugins=grpc:.
```

### 5. 服务端代码

```go
package main

import (
    "context"
    "fmt"
    "google.golang.org/grpc"
    "grpc_demo/hello"
    "net"
)

type Server struct {
}


func (s *Server)  SayHello(ctx context.Context,request *hello.HelloRequest)(*hello.HelloReply,error){
    return &hello.HelloReply{Message:"Hello "+request.Name},nil
}

func main()  {
    g := grpc.NewServer()
    s := Server{}
    hello.RegisterGreeterServer(g,&s)
    lis, err := net.Listen("tcp", fmt.Sprintf(":8080"))
    if err != nil {
        panic("failed to listen: "+err.Error())
    }
    g.Serve(lis)
}
```

### 6. 客户端

```go
package main

import (
    "context"
    "fmt"
    "google.golang.org/grpc"
    "grpc_demo/proto"
)

func main()  {
    conn,err := grpc.Dial("127.0.0.1:8080",grpc.WithInsecure())
    if err!=nil{
        panic(err)
    }
    defer conn.Close()
    c := hello.NewGreeterClient(conn)
    r,err := c.SayHello(context.Background(),&hello.HelloRequest{Name:"bobby"})
    if err!=nil{
        panic(err)
    }
    fmt.Println(r.Message)
}
```



## 3、protobuf和json的直观对比

编写proto文件 helloworld.proto

```protobuf
syntax = "proto3";
option go_package = "./;proto";  //新版本中需要加 /
message HelloRequest{
  string name = 1; //1是编号不是值
  int32 age = 2;
  repeated string courses = 3; //repeated 表示是一个切片，可以重复的值。
}
```

在控制台 cd 到 helloworld.proto 文件目录，运行下面命令 

```shell
protoc -I . helloworld.proto --go_out=plugins=grpc:. 
```

来生成go文件。

实际使用过程中，改用了 protoc --go_out=.  ./helloworld.proto 也可以生成成功，也就是使用了命令

```shell
protoc --go_out=output_directory input_directory/file.proto
```

新建proto编码测试文件，会发现protobuf协议打印出来的字符很少，比json占用的空间少，提高了传输效率。

```go
package main

import (
	hello "OldPackageTest/helloworld/proto"
	"encoding/json"
	"fmt"
)
import "github.com/golang/protobuf/proto"

type Hello struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Courses []string `json:"courses"`
}

func main() {
	req := hello.HelloRequest{
		Name:    "Tom",
		Age:     18,
		Courses: []string{"go", "gin", "微服务"},
	}
	rsp, _ := proto.Marshal(&req)
	fmt.Println(string(rsp))

	jsonStruct := Hello{Name: "Tom", Age: 18, Courses: []string{"go", "gin", "微服务"}}
	jsonRsp, _ := json.Marshal(jsonStruct)
	fmt.Println(string(jsonRsp))
}
```

反序列化测试

```go
package main

import (
	hello "OldPackageTest/helloworld/proto"
	"fmt"
)
import "github.com/golang/protobuf/proto"

type Hello struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Courses []string `json:"courses"`
}

func main() {
	req := hello.HelloRequest{
		Name:    "Tom",
		Age:     18,
		Courses: []string{"go", "gin", "微服务"},
	}
	rsp, _ := proto.Marshal(&req)
	fmt.Println(string(rsp))

	newReq := hello.HelloRequest{}
	_ = proto.Unmarshal(rsp, &newReq)
	fmt.Println(newReq.Name, newReq.Age, newReq.Courses)
}

```

