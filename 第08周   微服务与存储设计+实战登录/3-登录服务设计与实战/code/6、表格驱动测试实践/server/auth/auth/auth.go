package auth

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/dao"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service implements auth service.
type Service struct {
	OpenIDResolver OpenIDResolver //接口不用写*
	Mongo          *dao.Mongo
	Logger         *zap.Logger
}

// OpenIDResolver resolves an authorization code
// to an open id.
type OpenIDResolver interface {
	Resolve(code string) (string, error)
}

// Login logs a user in.
func (s *Service) Login(c context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	s.Logger.Info("received code", zap.String("code", req.Code))
	openID, err := s.OpenIDResolver.Resolve(req.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "connot resolve openid: %v", err)
	}
	accountID, err := s.Mongo.ResolveAccountID(c, openID)
	if err != nil { //小程序很难同时登陆，很难出索引错误
		//出现了很难出现的错误，记日志
		s.Logger.Error("received resolve accountID", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	return &authpb.LoginResponse{
		AccessToken: "token for account id:" + accountID,
		ExpiresIn:   7200,
	}, nil
}
