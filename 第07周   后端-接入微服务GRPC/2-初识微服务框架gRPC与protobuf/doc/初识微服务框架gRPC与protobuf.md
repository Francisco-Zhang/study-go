## 1、 GRPC的作用

基于http/2协议，与http/1协议相比

- 外层协议相同
- 二进制数据流，高效的传输 （不需要全部数据都准备好了再进行传输）
- 多路复用。 多个请求共用一个链接。http1 也可以共用，但是做不到 请求1 和 请求2 同时共用。
- 安全性提升

方法采用 POST

路径  /Service/Method 	  eg: 	/TripService/GetTrip



主要优点：

- 高效的数据传输
- 语言无关的领域模型定义



其他DSL/IDL

- Thrift
- Swagger
- Goa



## 2、ProtoBuf编译器的安装

GitHub release 中安装相应操作系统的最新版本

解压缩后，把执行文件的路径设置到环境变量里去。然后终端输入 protoc 验证



安装针对go语言进行解析到安装包

https://github.com/grpc-ecosystem/grpc-gateway

```sh
 go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```



## 3、ProtoBuf的使用

### vscode 安装插件 vscode-proto3

```protobuf
//sp3:语法提示
syntax = "proto3";
package coolcar;
option go_package="coolcar/proto/gen/go;trippb";

message Trip{
    string start =1;
    string end =2;
    int64 duration_sec=3;  //网络传输用_,转换成go后会替换为驼峰命名
    int64 fee_cent=4;
}
```

### 生成go代码

```shell
cd proto/;
mkdir -p gen/go

protoc -I=. --go_out=paths=source_relative:gen/go trip.proto
```

### 使用

```go
import (
	trippb "coolcar/proto/gen/go"
	"fmt"
)
func main() {
	trip:=trippb.Trip{}
	fmt.Println("hello world")
}
```

### 转换和解析

```go
package main

import (
	trippb "coolcar/proto/gen/go"
	"fmt"
	"google.golang.org/protobuf/proto"
)

func main() {
	trip := trippb.Trip{
		Start:       "abc",
		End:         "def",
		DurationSec: 3600,
		FeeCent:     1000,
	}
	fmt.Println(&trip)
	b, err := proto.Marshal(&trip) //用地址，防止赋值的时候使用私有变量
	if err != nil {
		panic(err)
	}
	fmt.Printf("%X\n", b) //以16进制打印字节流

	var trip2 trippb.Trip
	err = proto.Unmarshal(b, &trip2)
	if err != nil {
		panic(err)
	}
	fmt.Println(&trip2)
  
  b, err = json.Marshal(&trip2)  //可以转换成json，原因是结构体内定义了json标签
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)
}
```

## 4、复合类型和枚举类型

在proto文件中使用 shift+space 快捷键触发 trigger suggestion

proto中，bytes类型很适合传输图片等文件流。

```protobuf
syntax = "proto3";
package coolcar;
option go_package="coolcar/proto/gen/go;trippb";

message Location{
    double latitude=1;
    double longitude=2;
}
enum TripStatus{
    TS_NOT_SPECIFIED=0;
    NOT_STARTED=1;
    IN_PROGRESS=2;
    FINISHED=3;
    PAID=4;
}

message Trip{
    string start =1;
    Location start_pos=5;   //如果已经上线，旧版本的序号已经确定。新版本序号只能只能增加，不能占用修改。
    repeated Location path_locations=7; //repeated 变量一般加 s,表示复数。
    string end =2;
    Location end_pos=6;
    int64 duration_sec=3;
    int64 fee_cent=4;
    TripStatus status=8;
}
```

```go
package main

import (
	trippb "coolcar/proto/gen/go"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
)

func main() {
	trip := trippb.Trip{
		Start:       "abc",
		End:         "def",
		DurationSec: 3600,
		FeeCent:     1000,
		StartPos: &trippb.Location{
			Latitude:  30,
			Longitude: 120,
		},
		EndPos: &trippb.Location{
			Latitude:  35,
			Longitude: 115,
		},
		PathLocations: []*trippb.Location{
			{
				Latitude:  31,
				Longitude: 119,
			},
			{
				Latitude:  32,
				Longitude: 118,
			},
		},
		Status: trippb.TripStatus_FINISHED,
	}
	fmt.Println(&trip)
	b, err := proto.Marshal(&trip) //用地址，防止赋值的时候使用私有变量
	if err != nil {
		panic(err)
	}
	fmt.Printf("%X\n", b) //以16进制打印字节流

	var trip2 trippb.Trip
	err = proto.Unmarshal(b, &trip2)
	if err != nil {
		panic(err)
	}
	fmt.Println(&trip2)

	b, err = json.Marshal(&trip2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)
}
```

## 5、ProtoBuf字段的可选性

在上线新系统之后，新老系统的交互怎么解决？

新系统加了字段，老系统没有加，新系统向老系统传输数据，老系统proto在解析的时候会自动忽略。

反之老系统向新系统传输数据，新系统就需要使用可选字段，值就是0。

无法区分0是赋值后的0还是默认值0。



**如果某字段传输的数据为0，则不会传输该字段。**

```go
import (
	trippb "coolcar/proto/gen/go"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
)

func main() {
	trip := trippb.Trip{
		Start:       "abc",
		End:         "def",
		DurationSec: 0,
		FeeCent:     1000,
	}
	fmt.Println(&trip)
	b, err := proto.Marshal(&trip) //用地址，防止赋值的时候使用私有变量
	if err != nil {
		panic(err)
	}
	fmt.Printf("%X\n", b) //以16进制打印字节流

	var trip2 trippb.Trip
	err = proto.Unmarshal(b, &trip2)
	if err != nil {
		panic(err)
	}
	fmt.Println(&trip2)

	b, err = json.Marshal(&trip2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)
}
//打印数据结果中没有DurationSec，这就是tag中 omitempty 的含义：
start:"abc" end:"def" fee_cent:1000
0A03616263120364656620E807
start:"abc" end:"def" fee_cent:1000
{"start":"abc","end":"def","fee_cent":1000}
```

如果实在要区分呢，**需要新增字段** ,bool has_duration_sec=9;

```protobuf
message Trip{
    string start =1;
    Location start_pos=5;   //如果已经上线，旧版本的序号已经确定。新版本序号只能只能增加，不能占用修改。
    repeated Location path_locations=7; //repeated 变量一般加 s,表示复数。
    string end =2;
    Location end_pos=6;
    int64 duration_sec=3;
    int64 fee_cent=4;
    TripStatus status=8;
    bool has_duration_sec=9;
}
```

**这样每次新增字段会很麻烦，所以要巧用默认值。**

新系统新增字段  bool isPromotionTrip=9; 这样老系统不传，默认也是false，逻辑上就是对的。
