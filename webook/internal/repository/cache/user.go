package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"geek-basic-go/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

// UserCache
// 引入专门的UserCache:
// 1. 屏蔽过期时间设置问题：使用UserCache的人不再需要关心过期时间问题
// 2. 屏蔽key的结构，缓存调用者不用知道缓存里边这个key的结构是怎么组成的
// 3. 屏蔽序列化与反序列协议
type UserCache struct {
	// Cmdable有很多不同的实现，单机的，集群的
	cmd        redis.Cmdable
	expiration time.Duration
}

func (c *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := c.key(id)
	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	// 反序列化：将字节切片转换为结构体
	err = json.Unmarshal([]byte(data), &u)
	/*if err != nil {
		return domain.User{}, err
	}*/
	return u, err
}

func (c *UserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func (c *UserCache) Set(ctx context.Context, du domain.User) error {
	key := c.key(du.Id)
	// 序列化：将结构体转换为字节切片[]byte
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func NewUserCache(cmd redis.Cmdable) *UserCache {
	return &UserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

type UserCacheV1 struct {
	//面向了具体实现, 不推荐，推荐使用接口编程
	client *redis.Client
}

// NewUserCacheV1
// 一定不要自己去初始化程序需要的东西，让外边传进来
func NewUserCacheV1(addr string) *UserCache {
	cmd := redis.NewClient(&redis.Options{Addr: addr})
	return &UserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
