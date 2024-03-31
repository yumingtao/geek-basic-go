package ioc

import (
	"geek-basic-go/webook/internal/service/sms"
	"geek-basic-go/webook/internal/service/sms/localsms"
	"geek-basic-go/webook/internal/service/sms/tencent"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
)

func InitSmsService() sms.Service {
	// 如何使用装饰器
	//return ratelimit.NewLimitSmsService(localsms.NewService(), limiter.NewRedisSlidingWindowLimiter())
	//return opentelemetry.NewOtelSmsService(localsms.NewService(), NewTracer())
	return localsms.NewService()
	// 此处可以换成不同的实现
	//return InitTencentSmsService()
}

func InitTencentSmsService() sms.Service {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("找不到腾讯SMS的secret id")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("找不到腾讯SMS的secret key")
	}
	c, err := tencentSms.NewClient(common.NewCredential(secretId, secretKey), "ap-nanjing",
		profile.NewClientProfile())
	if err != nil {
		panic(err)
	}

	return tencent.NewService(c, "123456789", "yumingtao", nil)
}
