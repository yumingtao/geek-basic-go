package ioc

import (
	"geek-basic-go/webook/internal/service/oauth2/wechat"
	"os"
)

func InitWechatService() wechat.Service {
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("WECHAT_APP_ID environment variable not set")
	}
	return wechat.NewService(appId)
}
