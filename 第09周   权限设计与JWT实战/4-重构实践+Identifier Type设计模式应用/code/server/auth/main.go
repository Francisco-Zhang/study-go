package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/auth/wechat"
	"coolcar/auth/dao"
	"coolcar/auth/token"
	"coolcar/shared/server"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	c := context.Background()
	mc, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/coolcar?readPreference=primary&ssl=false"))
	if err != nil {
		log.Fatal("cannot connect mongodb", zap.Error(err))
	}

	pkFile, err := os.Open("private.key")
	if err != nil {
		log.Fatal("cannot open private key", zap.Error(err))
	}
	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		log.Fatal("cannot read private key", zap.Error(err))
	}
	privkey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		log.Fatal("cannot parse private key", zap.Error(err))
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:   "auth",
		Addr:   ":8081",
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
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
		},
	}))

}
