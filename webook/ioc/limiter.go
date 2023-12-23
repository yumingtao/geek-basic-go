package ioc

import (
	"geek-basic-go/webook/pkg/limiter"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitLimiter(cmd redis.Cmdable) limiter.Limiter {
	return limiter.NewRedisSlidingWindowLimiter(cmd, time.Second, 100)
}
