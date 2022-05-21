package main

import (
	trippb "coolcar/proto/gen/go"
	trip "coolcar/tripservice"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
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
