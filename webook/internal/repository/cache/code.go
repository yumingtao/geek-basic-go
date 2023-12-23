package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"time"
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type GoCacheCodeCache struct {
	cache *cache.Cache
}

type RedisCodeCache struct {
	cmd redis.Cmdable
}

const (
	CodeExpiration = time.Minute * 10
	KeyCnt         = 3
)

var (
	//go:embed lua/set_code.lua
	lusSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode        string
	ErrCodeSendTooMany   = errors.New("发送太频繁")
	ErrVerifyCodeTooMany = errors.New("验证太频繁")
)

func NewGoCacheCodeCache() CodeCache {
	return &GoCacheCodeCache{
		// 过期时间CodeExpirationSeconds=600s，cleanup interval 设置为-1，手动清理
		cache: cache.New(CodeExpiration, -1*time.Second),
	}
}

func NewRedisCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}
}

func (c *GoCacheCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	key := key(biz, phone)
	cntKey := cntKey(biz, phone)
	now := time.Now()
	_, expiration, found := c.cache.GetWithExpiration(key)
	// 注意Set方法会替换掉缓存里已经存在的值
	// 验证码存在
	if found {
		// 过期时间是expiration是zero value
		if expiration.IsZero() {
			return errors.New("验证码存在，但是没有过期时间")
		}
		// 过期时间大于9分钟，验证码发送太过频繁
		if expiration.Sub(now).Seconds() >= 540 {
			return ErrCodeSendTooMany
		}
	}

	// 下面两种情况下才Set
	// 1. 没找到，或是已经过期了
	// 2. 过期时间小于9分钟
	c.cache.Set(key, code, CodeExpiration)
	// 设置验证码时，不用判断验证次数是否在缓存中存在，直接替换
	c.cache.Set(cntKey, KeyCnt, CodeExpiration)

	return nil
}

func (c *GoCacheCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	cntKey := cntKey(biz, phone)
	cachedCnt, found := c.cache.Get(cntKey)
	if !found {
		// 没找到，说明已经过期或是
		return false, nil
	} else {
		// 找到，判断次数
		cnt, ok := cachedCnt.(int)
		if !ok {
			return false, nil
		}
		if cnt <= 0 {
			return false, ErrVerifyCodeTooMany
		}
	}

	key := key(biz, phone)
	cachedCode, found := c.cache.Get(key)
	if !found {
		return false, nil
	} else {
		if cachedCode == code {
			c.cache.Set(cntKey, 0, CodeExpiration)
			return true, nil
		} else {
			err := c.cache.Decrement(cntKey, 1)
			return false, err
		}
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.cmd.Eval(ctx, lusSetCode, []string{key(biz, phone)}, code).Int()
	if err != nil {
		// 调用redis出现问题
		return err
	}

	switch res {
	case -2:
		return errors.New("验证码存在，但是没有过期时间")
	case -1:
		return ErrCodeSendTooMany
	default:
		return nil
	}
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{key(biz, phone)}, code).Int()
	if err != nil {
		// 调用redis出现问题
		return false, err
	}

	switch res {
	case -1:
		return false, ErrVerifyCodeTooMany
	case -2:
		return false, nil
	default:
		return true, nil
	}
}

func key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func cntKey(biz, phone string) string {
	return key(biz, phone) + ":cnt"
}
