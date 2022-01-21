## 1、什么是rpc? rpc 开发的挑战是什么？

### 什么是rpc

1. RPC（Remote Procedure Call）远程过程调用，简单理解是一个节点请求另一个节点提供的服务。
2. 对应rpc的是本地过程调用，函数调用是最常见的本地过程调用。
3. 将本地过程调用变成远程过程调用会面临各种问题。



### 远程过程面临的问题

1. 原本的本地函数放到另一个服务器上去运行。但是引入了很多新问题
2. Call 的 id 映射
3. 序列化 和 反序列化
4. 网路传输



在远程调用时，我们需要执行的函数体是在远程的机器上的，也就是说，add是在另一个进程中执行的。这就带来了几个新问题：

1. Call ID映射。我们怎么告诉远程机器我们要调用add，而不是sub或者Foo呢？在本地调用中，函数体是直接通过函数指针来指定的，我们调用add，编译器就自动帮我们调用它相应的函数指针。但是在远程调用中，函数指针是不行的，因为两个进程的地址空间是完全不一样的。所以，在RPC中，所有的函数都必须有自己的一个ID。这个ID在所有进程中都是唯一确定的。客户端在做远程过程调用时，必须附上这个ID。然后我们还需要在客户端和服务端分别维护一个 {函数 <--> Call ID} 的对应表。两者的表不一定需要完全相同，但相同的函数对应的Call ID必须相同。当客户端需要进行远程调用时，它就查一下这个表，找出相应的Call ID，然后把它传给服务端，服务端也通过查表，来确定客户端需要调用的函数，然后执行相应函数的代码。
2. 序列化和反序列化。客户端怎么把参数值传给远程的函数呢？在本地调用中，我们只需要把参数压到栈里，然后让函数自己去栈里读就行。但是在远程过程调用时，客户端跟服务端是不同的进程，不能通过内存来传递参数。甚至有时候客户端和服务端使用的都不是同一种语言（比如服务端用C++，客户端用Java或者Python）。这时候就需要客户端把参数先转成一个字节流，传给服务端后，再把字节流转成自己能读取的格式。这个过程叫序列化和反序列化。同理，从服务端返回的值也需要序列化反序列化的过程。
3. 网络传输。远程调用往往用在网络上，客户端和服务端是通过网络连接的。所有的数据都需要通过网络传输，因此就需要有一个网络传输层。网络传输层需要把Call ID和序列化后的参数字节流传给服务端，然后再把序列化后的调用结果传回客户端。只要能完成这两者的，都可以作为传输层使用。因此，它所使用的协议其实是不限的，能完成传输就行。尽管大部分RPC框架都使用TCP协议，但其实UDP也可以，而gRPC干脆就用了HTTP2。Java的Netty也属于这层的东西。

**解决了上面三个机制，就能实现RPC了，具体过程如下：**

client端解决的问题：

```plain
1. 将这个调用映射为Call ID。这里假设用最简单的字符串当Call ID的方法
2. 将Call ID，a和b序列化。可以直接将它们的值以二进制形式打包
3. 把2中得到的数据包发送给ServerAddr，这需要使用网络传输层
4. 等待服务器返回结果
4. 如果服务器调用成功，那么就将结果反序列化，并赋给total
```

server端解决的问题

```plain
1. 在本地维护一个Call ID到函数指针的映射call_id_map，可以用dict完成
2. 等待请求，包括多线程的并发处理能力
3. 得到一个请求后，将其数据包反序列化，得到Call ID
4. 通过在call_id_map中查找，得到相应的函数指针
5. 将a和rb反序列化后，在本地调用add函数，得到结果
6. 将结果序列化后通过网络返回给Client
```

在上面的整个流程中，估计有部分同学看到了熟悉的计算机网络的流程和web服务器的定义。

所以要实现一个RPC框架，其实只需要按以上流程实现就基本完成了。



其中：

- Call ID映射可以直接使用函数字符串，也可以使用整数ID。映射表一般就是一个哈希表。
- 序列化反序列化可以自己写，也可以使用Protobuf或者FlatBuffers之类的。

- 网络传输库可以自己写socket，或者用asio，ZeroMQ，Netty之类。



实际上真正的开发过程中，除了上面的基本功能以外还需要更多的细节：网络错误、流量控制、超时和重试等。

http2.0解决了http1.0一次请求响应就连接断开的问题，也解决了自定义协议通用性不高的问题。grpc 选择的传输协议就是http2.0。



## 2、通过 http 完成 add 服务端功能

```go
func main() {
	//http://127.0.0.1:8000/add?a=3&b=4
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm() //解析参数
		fmt.Println("path", r.URL.Path)
		a, _ := strconv.Atoi(r.Form["a"][0])
		b, _ := strconv.Atoi(r.Form["b"][0])
		w.Header().Set("Content-Type", "application/json")
		jData, _ := json.Marshal(map[string]int{"data": a + b})
		_, _ = w.Write(jData)
	})
	_ = http.ListenAndServe(":8000", nil)
}
```

## 3、通过 http 完成 add 客户端的功能

```go
import (
	"encoding/json"
	"fmt"
	"github.com/kirinlabs/HttpRequest"
)

type ResponseData struct {
	Data int `json:"data"`
}

func Add(a, b int) int {
	req := HttpRequest.NewRequest()
	res, _ := req.Get(fmt.Sprintf("http://localhost:8000/%s?a=%d&b=%d", "add", a, b))
	body, _ := res.Body()
	//fmt.Println(string(body))
	rspData := ResponseData{}
	json.Unmarshal(body, &rspData)
	return rspData.Data
}

func main() {
	fmt.Println(Add(3, 4))
}
```



## 4、rpc 架构分析

### **RPC框架的组成**

一个基本RPC框架由4部分组成，分别是：客户端、客户端存根、服务端、服务端存根

- 客户端（Client）：服务调用的发起方，也称为消费者。

- 客户端存根（ClientStub）：该程序运行在客户端所在的计算机上，主要用来存储 要调用服务器地址，对客户端发送的数据进行序列化、建立网络连接之后发送数据包给Server端的存根程序。接收Server端存根程序响应的消息，对消息进行反序列化。

- 服务端（Server）:提供客户端想要调用的函数

- 服务端存根（ServerStub）:接收客户端的消息并反序列化，对server端响应的消息进行序列化并响应给Client端

总体来讲客户端和服务端的Stub在底层帮助我们实现了Call ID的映射、数据的序列化和反序列化、网络传输，这样RPC客户端和RPC服务端 只需要专注于服务端和客户端的业务逻辑层。

![1](img/1.png)