package ioc

import (
	rlock "github.com/gotomicro/redis-lock"
	"github.com/redis/go-redis/v9"
)

func InitRLockClient(client redis.Cmdable) *rlock.Client {
	return rlock.NewClient(client)
}
