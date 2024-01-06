package startup

import (
	"geek-basic-go/webook/internal/service/oauth2/wechat"
	"geek-basic-go/webook/pkg/logger"
)

func InitWechatService(l logger.LoggerV1) wechat.Service {
	return wechat.NewService("appId", "appSecret", l)
}
