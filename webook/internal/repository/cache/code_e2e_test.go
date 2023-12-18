package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGoCacheCodeCache_Set_e2e(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	testCases := []struct {
		name      string
		before    func(t *testing.T)
		after     func(t *testing.T)
		ctx       context.Context
		biz       string
		phone     string
		code      string
		wantedErr error
	}{
		{
			name: "设置成功",
			before: func(t *testing.T) {
				// 准备数据
			},
			after: func(t *testing.T) {
				// 验证数据
				// 需要验证验证码存在了redis里
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13813382231"
				dur, err := rdb.TTL(ctx, key).Result()
				assert.True(t, dur > time.Minute*9+time.Second+50)
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			ctx:       context.Background(),
			biz:       "login",
			phone:     "13813382231",
			code:      "123456",
			wantedErr: nil,
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				// 准备数据
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13813382231"
				err := rdb.Set(ctx, key, "654321", time.Minute*9+time.Second+50).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据
				// 需要验证验证码存在了redis里
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13813382231"
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "654321", code)
			},
			ctx:       context.Background(),
			biz:       "login",
			phone:     "13813382231",
			code:      "123456",
			wantedErr: ErrCodeSendTooMany,
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				// 准备数据
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13813382231"
				// redis中有key但是没有过期时间，web/user.go中的SendSmsLoginCode代码会进入default系统错误分支
				err := rdb.Set(ctx, key, "654321", 0).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据
				// 需要验证验证码存在了redis里
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13813382231"
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "654321", code)
			},
			ctx:       context.Background(),
			biz:       "login",
			phone:     "13813382231",
			code:      "123456",
			wantedErr: errors.New("验证码存在，但是没有过期时间"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			c := NewRedisCodeCache(rdb)
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantedErr, err)
		})
	}
}
