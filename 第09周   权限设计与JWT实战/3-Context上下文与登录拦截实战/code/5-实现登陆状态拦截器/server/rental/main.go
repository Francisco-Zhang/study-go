package main

import (
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip"
	"coolcar/shared/auth"
	"log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := newZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		//输出完日志后，程序自动退出
		logger.Fatal("cannot  listen:", zap.Error(err))
	}

	in, err := auth.Interceptor("../shared/auth/public.key")
	if err != nil {
		//输出完日志后，程序自动退出
		logger.Fatal("cannot  create Interceptor:", zap.Error(err))
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(in))

	//服务里是 *Service，这里就要取地址
	rentalpb.RegisterTripServiceServer(s, &trip.Service{
		Logger: logger,
	})
	err = s.Serve(lis)
	logger.Fatal("cannot serve:", zap.Error(err))
}

func newZapLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey = ""
	return cfg.Build()
}
