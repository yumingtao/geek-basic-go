package wechat

import (
	"context"
	"fmt"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/url"
)

var redirectURL = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

const (
	authURLPattern = `https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`
)

type Service interface {
	AuthUrl(ctx context.Context) (string, error)
}

type WechatService struct {
	appId string
}

func NewService(appId string) Service {
	return WechatService{
		appId: appId,
	}
}

func (w WechatService) AuthUrl(ctx context.Context) (string, error) {
	state := uuid.New()
	return fmt.Sprintf(authURLPattern, w.appId, redirectURL, state), nil
}
