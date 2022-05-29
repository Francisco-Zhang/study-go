package wechat

import (
	"fmt"

	"github.com/medivhzhan/weapp/v2"
)

type Service struct {
	AppID     string
	AppSecret string
}

func (s *Service) Resolve(code string) (string, error) {
	resp, err := weapp.Login(s.AppID, s.AppSecret, code)
	if err != nil {
		return "nil", fmt.Errorf("weapp login: %v", err)
	}
	if err := resp.GetResponseError(); err != nil {
		return "nil", fmt.Errorf("weapp response error: %v", err)
	}
	return resp.OpenID, nil
}
