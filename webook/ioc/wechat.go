package ioc

import (
	"geek-basic-go/webook/internal/service/oauth2/wechat"
	"geek-basic-go/webook/pkg/logger"
	"os"
)

func InitWechatService(l logger.LoggerV1) wechat.Service {
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("WECHAT_APP_ID environment variable not set")
	}
	appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("WECHAT_APP_SECRET environment variable not set")
	}
	return wechat.NewService(appId, appSecret, l)
}
