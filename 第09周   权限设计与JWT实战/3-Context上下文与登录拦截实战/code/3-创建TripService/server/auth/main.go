package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/auth/wechat"
	"coolcar/auth/dao"
	"coolcar/auth/token"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	c := context.Background()
	mc, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/coolcar?readPreference=primary&ssl=false"))
	if err != nil {
		log.Fatalf("cannot connect mongodb", zap.Error(err))
	}

	pkFile, err := os.Open("private.key")
	if err != nil {
		log.Fatalf("cannot open private key", zap.Error(err))
	}
	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		log.Fatalf("cannot read private key", zap.Error(err))
	}
	privkey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		log.Fatalf("cannot parse private key", zap.Error(err))
	}
	s := grpc.NewServer()
	//服务里是 *Service，这里就要取地址
	authpb.RegisterAuthServiceServer(s, &auth.Service{
		OpenIDResolver: &wechat.Service{
			AppID:     "wx692887841c1ef470",
			AppSecret: "571ff2da753b8ca0f2985714dd04b039",
		}, // *Service 才有login，此处取地址
		Mongo:          dao.NewMongo(mc.Database("coolcar")),
		Logger:         logger,
		TokenExpire:    2 * time.Hour,
		TokenGenerator: token.NewJWTTokenGen("car/auth", privkey),
	})
	err = s.Serve(lis)
	logger.Fatal("cannot serve:", zap.Error(err))
}

func newZapLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey = ""
	return cfg.Build()
}
