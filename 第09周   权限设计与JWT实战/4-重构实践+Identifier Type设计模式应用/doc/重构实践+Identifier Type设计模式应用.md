## 1、重构

### shared/server/grpc.go

```go
package server

import (
	"coolcar/shared/auth"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// GRPCConfig defines a grpc server.
type GRPCConfig struct {
	Name              string
	Addr              string
	AuthPublicKeyFile string
	RegisterFunc      func(*grpc.Server)
	Logger            *zap.Logger
}

// RunGRPCServer runs a grpc server.
func RunGRPCServer(c *GRPCConfig) error {
	nameField := zap.String("name", c.Name)
	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		c.Logger.Fatal("cannot listen", nameField, zap.Error(err))
	}

	var opts []grpc.ServerOption
	if c.AuthPublicKeyFile != "" {
		in, err := auth.Interceptor(c.AuthPublicKeyFile)
		if err != nil {
			c.Logger.Fatal("cannot create auth interceptor", nameField, zap.Error(err))
		}
		opts = append(opts, grpc.UnaryInterceptor(in))
	}

	s := grpc.NewServer(opts...)
	c.RegisterFunc(s)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	c.Logger.Info("server started", nameField, zap.String("addr", c.Addr))
	return s.Serve(lis)
}
```

### shared/server/logger.go

```go
package server

import "go.uber.org/zap"

// NewZapLogger creates a new zap logger
func NewZapLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey = ""
	return cfg.Build()
}
```

### rental/main.go

```go
package main

import (
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip"
	"coolcar/shared/server"
	"log"

	"google.golang.org/grpc"
)

func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:              "rental",
		Addr:              ":8082",
		AuthPublicKeyFile: "../shared/auth/public.key",
		Logger:            logger,
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				Logger: logger,
			})
		},
	}))
}
```

## 2、应用Identifier Type设计模式

这样设计可以防止出错，比如 tripId,userId都是string类型，互相赋值，编译器没法检查出来。

```go
package id

// AccountID defines account id object.
type AccountID string

func (a AccountID) String() string {
	return string(a)
}


func (i *interceptor) HandleReq(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	tkn, err := tokenFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "")
	}
	aid, err := i.verifier.Verify(tkn)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token not valid: %v", err)
	}

	//其他领域也需要使用AccountID，所以方法大写
	return handler(ContextWithAccountID(ctx, id.AccountID(aid)), req)
}


func ContextWithAccountID(c context.Context, aid id.AccountID) context.Context {
	return context.WithValue(c, accountIDKey{}, aid)
}


func AccountIDFromContext(c context.Context) (id.AccountID, error) {
	v := c.Value(accountIDKey{})
	aid, ok := v.(id.AccountID)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "")
	}
	return aid, nil
}

```

