package main

import (
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/auth/wechat"
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
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		//输出完日志后，程序自动退出
		logger.Fatal("cannot  listen:", zap.Error(err))
	}
	s := grpc.NewServer()
	//服务里是 *Service，这里就要取地址
	authpb.RegisterAuthServiceServer(s, &auth.Service{
		OpenIDResolver: &wechat.Service{
			AppID:     "wx692887841c1ef470",
			AppSecret: "571ff2da753b8ca0f2985714dd04b039",
		}, // *Service 才有login，此处取地址
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
