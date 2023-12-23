package repository

import (
	"context"
	"database/sql"
	"errors"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository/cache"
	cachemocks "geek-basic-go/webook/internal/repository/cache/mocks"
	"geek-basic-go/webook/internal/repository/dao"
	daomocks "geek-basic-go/webook/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	nowMs := time.Now().UnixMilli()
	now := time.UnixMilli(nowMs)
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao)
		ctx        context.Context
		uid        int64
		wantedUser domain.User
		wantedErr  error
	}{
		{
			name: "查找成功，缓存未命中",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				uid := int64(123)
				c := cachemocks.NewMockUserCache(ctrl)
				d := daomocks.NewMockUserDao(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).Return(dao.User{
					Id:       uid,
					NickName: "456",
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Password:        "123456",
					BirthDate:       "1983-03-30",
					PersonalProfile: "我是一个好人",
					Phone: sql.NullString{
						String: "123456778",
						Valid:  true,
					},
					CreatedAt: nowMs,
					UAt:       nowMs,
				}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:              123,
					NickName:        "456",
					Email:           "123@qq.com",
					Password:        "123456",
					BirthDate:       "1983-03-30",
					PersonalProfile: "我是一个好人",
					Phone:           "123456778",
					Ctime:           now,
				}).Return(nil)
				return c, d
			},
			ctx: context.Background(),
			uid: 123,
			wantedUser: domain.User{
				Id:              123,
				NickName:        "456",
				Email:           "123@qq.com",
				Password:        "123456",
				BirthDate:       "1983-03-30",
				PersonalProfile: "我是一个好人",
				Phone:           "123456778",
				Ctime:           now,
			},
			wantedErr: nil,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				uid := int64(123)
				c := cachemocks.NewMockUserCache(ctrl)
				d := daomocks.NewMockUserDao(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).Return(domain.User{
					Id:              123,
					NickName:        "456",
					Email:           "123@qq.com",
					Password:        "123456",
					BirthDate:       "1983-03-30",
					PersonalProfile: "我是一个好人",
					Phone:           "123456778",
					Ctime:           now,
				}, nil)
				return c, d
			},
			ctx: context.Background(),
			uid: 123,
			wantedUser: domain.User{
				Id:              123,
				NickName:        "456",
				Email:           "123@qq.com",
				Password:        "123456",
				BirthDate:       "1983-03-30",
				PersonalProfile: "我是一个好人",
				Phone:           "123456778",
				Ctime:           now,
			},
			wantedErr: nil,
		},
		{
			name: "未找到用户",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				uid := int64(123)
				c := cachemocks.NewMockUserCache(ctrl)
				d := daomocks.NewMockUserDao(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).Return(dao.User{}, dao.ErrRecordNotFound)
				return c, d
			},
			ctx:        context.Background(),
			uid:        123,
			wantedUser: domain.User{},
			wantedErr:  ErrUserNotFound,
		},
		{
			name: "回写缓存失败",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				uid := int64(123)
				c := cachemocks.NewMockUserCache(ctrl)
				d := daomocks.NewMockUserDao(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).Return(dao.User{
					Id:       uid,
					NickName: "456",
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Password:        "123456",
					BirthDate:       "1983-03-30",
					PersonalProfile: "我是一个好人",
					Phone: sql.NullString{
						String: "123456778",
						Valid:  true,
					},
					CreatedAt: nowMs,
					UAt:       nowMs,
				}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:              123,
					NickName:        "456",
					Email:           "123@qq.com",
					Password:        "123456",
					BirthDate:       "1983-03-30",
					PersonalProfile: "我是一个好人",
					Phone:           "123456778",
					Ctime:           now,
				}).Return(errors.New("redis错误"))
				return c, d
			},
			ctx: context.Background(),
			uid: 123,
			wantedUser: domain.User{
				Id:              123,
				NickName:        "456",
				Email:           "123@qq.com",
				Password:        "123456",
				BirthDate:       "1983-03-30",
				PersonalProfile: "我是一个好人",
				Phone:           "123456778",
				Ctime:           now,
			},
			wantedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userCache, userDao := tc.mock(ctrl)
			repo := NewCachedUserRepository(userDao, userCache)
			user, err := repo.FindById(tc.ctx, tc.uid)
			assert.Equal(t, tc.wantedErr, err)
			assert.Equal(t, tc.wantedUser, user)
		})
	}
}
