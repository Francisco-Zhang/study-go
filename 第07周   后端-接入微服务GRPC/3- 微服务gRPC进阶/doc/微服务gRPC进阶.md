## 1、GRPC服务器及客户端

proto带方法生成命令：protoc -I=. --go_out=plugins=grpc,paths=source_relative:gen/go trip.proto

##  2、REST vs RPC

![1](img/1.png)

## 3、GRPC Gateway的作用

## 4、GRPC Gateway的实现

### 需要提前安装gateway包

```shell
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

### proto目录新建 trip.yaml配置文件

```yaml
type: google.api.Service
config_version: 3 

http:
  rules:
  - selector: coolcar.TripService.GetTrip
    get: /trip/{id}
```

### 通过脚本生成go代码

新建 gen.sh文件

```sh
protoc -I=. --go_out=plugins=grpc,paths=source_relative:gen/go trip.proto;
protoc -I=. --grpc-gateway_out=paths=source_relative,grpc_api_configuration=trip.yaml:gen/go trip.proto;
```

### 启动http服务，配置远程调用

```go
import (
	"context"
	trippb "coolcar/proto/gen/go"
	trip "coolcar/tripservice"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.SetFlags(log.Lshortfile) //设置日志开头为文件位置。
	go startGRPCGateway()
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		//输出完日志后，程序自动退出
		log.Fatalf("failed to listen: v%", err)
	}
	s := grpc.NewServer()
	//服务里是 *Service，这里就要取地址
	trippb.RegisterTripServiceServer(s, &trip.Service{})

	log.Fatal(s.Serve(lis))
}

func startGRPCGateway() {
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel() //一旦调用，会断开对grpc服务的连接

	mux := runtime.NewServeMux()
	err := trippb.RegisterTripServiceHandlerFromEndpoint(c, mux, "localhost:8081",
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		log.Fatalf("cannot start grpc gateway: %v", err)
	}
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("cannot listen and serve: %v", err)
	}
}
```

通过浏览器 访问接口 http://localhost:8080/trip/123，就能看到gateway的返回值。
