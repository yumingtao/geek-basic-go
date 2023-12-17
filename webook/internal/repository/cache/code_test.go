package cache

import (
	"context"
	"errors"
	"fmt"
	"geek-basic-go/webook/internal/repository/cache/redismocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	keyFun := func(biz, phone string) string {
		return fmt.Sprintf("phone_code:%s::%s", biz, phone)
	}
	testCases := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) redis.Cmdable
		ctx       context.Context
		biz       string
		phone     string
		code      string
		wantedErr error
	}{
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				redisCache := redismocks.NewMockCmdable(ctrl)
				// mock Eval返回的cmd
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(0))
				//cmd := redis.NewCmdResult(int64(0), nil)
				redisCache.EXPECT().Eval(
					gomock.Any(),
					lusSetCode,
					[]string{keyFun("test", "12345678")},
					"123456").Return(cmd)
				return redisCache
			},
			ctx:       context.Background(),
			biz:       "test",
			phone:     "12345678",
			code:      "123456",
			wantedErr: nil,
		},
		{
			name: "redis返回error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				redisCache := redismocks.NewMockCmdable(ctrl)
				// mock Eval返回的cmd
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(errors.New("redis错误"))
				redisCache.EXPECT().Eval(
					gomock.Any(),
					lusSetCode,
					[]string{keyFun("test", "12345678")},
					"123456").Return(cmd)
				return redisCache
			},
			ctx:       context.Background(),
			biz:       "test",
			phone:     "12345678",
			code:      "123456",
			wantedErr: errors.New("redis错误"),
		},
		{
			name: "验证码存在，但是没有过期时间",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				redisCache := redismocks.NewMockCmdable(ctrl)
				// mock Eval返回的cmd
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(-2))
				redisCache.EXPECT().Eval(
					gomock.Any(),
					lusSetCode,
					[]string{keyFun("test", "12345678")},
					"123456").Return(cmd)
				return redisCache
			},
			ctx:       context.Background(),
			biz:       "test",
			phone:     "12345678",
			code:      "123456",
			wantedErr: errors.New("验证码存在，但是没有过期时间"),
		},
		{
			name: "发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				redisCache := redismocks.NewMockCmdable(ctrl)
				// mock Eval返回的cmd
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(-1))
				redisCache.EXPECT().Eval(
					gomock.Any(),
					lusSetCode,
					[]string{keyFun("test", "12345678")},
					"123456").Return(cmd)
				return redisCache
			},
			ctx:       context.Background(),
			biz:       "test",
			phone:     "12345678",
			code:      "123456",
			wantedErr: ErrCodeSendTooMany,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewRedisCodeCache(tc.mock(ctrl))
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantedErr, err)
		})
	}
}
