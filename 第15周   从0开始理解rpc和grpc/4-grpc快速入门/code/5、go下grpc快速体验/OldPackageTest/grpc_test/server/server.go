package main

import (
	"OldPackageTest/grpc_test/proto"
	"context"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
}

// SayHello 第一个参数一定要是context,ctx 主要解决协程超时
func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {

	return &proto.HelloReply{Message: "hello," + request.Name}, nil
}

func main() {
	g := grpc.NewServer()
	proto.RegisterGreeterServer(g, &Server{})
	lis, err := net.Listen("tcp", "0.0.0.0:8088")
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	err = g.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
