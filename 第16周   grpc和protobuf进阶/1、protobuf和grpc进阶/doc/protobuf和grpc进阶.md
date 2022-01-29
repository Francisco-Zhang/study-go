## 1、option go_package的作用

```protobuf
option go_package = "common/stream/v1";
```

最后面的v1，同时定义了路径和包名。当末尾路径和包名一致时使用这种写法更方便，可以不使用之前`;`的定义方法。

![1](img/1.png)

## 2、proto文件同步时的坑

服务端的proto文件和客户端的proto文件不一致，就有可能出现各种解析问题。message 内字段的编号，两边一定要一致。

## 3、proto文件中import另一个proto文件

### rpc方法没有参数，但是必须要有占位

```protobuf
syntax = "proto3";
option go_package = ".;proto";
service Greeter {
    rpc SayHello (HelloRequest) returns (HelloReply);
    rpc Ping(Empty) returns (Pong);
}

message HelloRequest {
    string url = 1;
    string name = 2;
}

message HelloReply {
    string message = 1;
}
message Empty{  //需要定义一个空的message
}
message Pong{
    string id = 1;
}
```

### 公用message提取

由于多个proto文件都需要Ping,但是同一个message在一个包下面只能有一个。所以需要放到一个单独的包内。然后在其他的proto文件内使用 import 来引入。

**base.proto**

```protobuf
syntax = "proto3";
option go_package = ".;proto";

message Empty{ 
}
message Pong{
    string id = 1;
}
```

**hello.proto**

```protobuf
syntax = "proto3";
import "base.proto";
option go_package = ".;proto";
service Greeter {
    rpc SayHello (HelloRequest) returns (HelloReply);
    rpc Ping(Empty) returns (Pong);
}

message HelloRequest {
    string url = 1;
    string name = 2;
}

message HelloReply {
    string message = 1;
}
```

### 内置Empty

```protobuf
syntax = "proto3";
import "base.proto";
import "google/protobuf/empty.proto";
option go_package = ".;proto";
service Greeter {
    rpc SayHello (HelloRequest) returns (HelloReply);
    rpc Ping(google.protobuf.Empty) returns (Pong);
}

message HelloRequest {
    string url = 1;
    string name = 2;
}

message HelloReply {
    string message = 1;
}

```

### 如何在go中使用import的proto定义的message

1、找到import的proto源码，找到对应的 go_package，通过 go_package 定义的路径来确定 生成go 代码中的包名。 



## 4、嵌套的message对象

```protobuf
message HelloReply {
    string message = 1;

    message Result {	//放到HelloReply里边是为了不想做成公用的，只在Reply内部使用。
        string name = 1;
        string url = 2;
    }

    repeated Result data = 2;
}

```

在go代码中实例化 Result : proto.HelloReply_Result{} ，Result 名字会被修改。



## 5、protobuf中的enum枚举类型

```protobuf
enum Gender{
    MALE = 0;
    FEMALE = 1;
}

message HelloRequest {
    string name = 1; 
    string url = 2;
    Gender g = 3;
}
```

```go
rsp, _ := client.SayHello(context.Background(), &proto_bak.HelloRequest{
		Name: "bobby",
		Url:  "https://imooc.com",
		G:    proto_bak.Gender_MALE,
	})
```

## 6、map 类型 和 timestamp类型

```protobuf

message HelloRequest {
    string name = 1; //姓名 相当于文档
    string url = 2;
    Gender g = 3;
    map<string, string> mp = 4; // 不容易写文档注释，所以不要太频繁使用。
    google.protobuf.Timestamp addTime = 5;
}
```

```go
rsp, _ := client.SayHello(context.Background(), &proto_bak.HelloRequest{
		Name: "bobby",
		Url:  "https://imooc.com",
		G:    proto_bak.Gender_MALE,
		Mp: map[string]string{
			"name":    "bobby",
			"company": "慕课网",
		},
		AddTime: timestamppb.New(time.Now()),
	})
```

