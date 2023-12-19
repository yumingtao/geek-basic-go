package ratelimit

import (
	"context"
	"errors"
	"geek-basic-go/webook/internal/service/sms"
	"geek-basic-go/webook/pkg/limiter"
)

var errLimited = errors.New("触发限流")

// LimitSmsServiceV1 is a type that decorates the sms.Service interface with rate limiting functionality.
// 组合的方式的装饰器
// 1. 用户可以直接访问Service，绕开装饰器本身 2. 可以只实现service的部分方法
type LimitSmsServiceV1 struct {
	// 被装饰的
	sms.Service
	limiter limiter.Limiter
	key     string
}

// LimitSmsService is a type that decorates the sms.Service interface with rate limiting functionality.
// 不使用组合
// 1. 可以有效阻止用户绕开装饰器 2. 必须实现Service的全部方法
type LimitSmsService struct {
	// 被装饰的
	svc     sms.Service
	limiter limiter.Limiter
	key     string
}

func NewLimitSmsService(svc sms.Service, limiter limiter.Limiter) *LimitSmsService {
	return &LimitSmsService{
		svc:     svc,
		limiter: limiter,
		key:     "sms-limiter",
	}
}

func (l *LimitSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limited, err := l.limiter.Limit(ctx, l.key)
	if err != nil {
		return err
	}
	if limited {
		return errLimited
	}
	return l.svc.Send(ctx, tplId, args, numbers...)
}
