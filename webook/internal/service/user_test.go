package service

import (
	"context"
	"errors"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository"
	repomocks "geek-basic-go/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordEncrypt(t *testing.T) {
	pwd := []byte("123456*hello")
	encrypted, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	assert.NoError(t, err)
	println(string(encrypted))
	err = bcrypt.CompareHashAndPassword(encrypted, pwd)
	require.NoError(t, err)
}

func TestUserServiceImpl_Login(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) repository.UserRepository
		ctx        context.Context
		email      string
		password   string
		wantedUser domain.User
		wantedErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email: "123@qq.com",
					// 应该是一个加密后正确的密码
					Password: "$2a$10$FgqaFMtlGQomZU9Pi/26MeFGTVOrg0bjDzjcu2c/F6ZBx.MQoEP/O",
					Phone:    "12345678",
				}, nil)
				return repo
			},
			//ctx:   context.Background(),
			email: "123@qq.com",
			// 用户输入的没有加密的密码
			password: "123456*hello",
			wantedUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$FgqaFMtlGQomZU9Pi/26MeFGTVOrg0bjDzjcu2c/F6ZBx.MQoEP/O",
				Phone:    "12345678",
			},
			wantedErr: nil,
		},
		{
			name: "用户未找到",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			//ctx:   context.Background(),
			email: "123@qq.com",
			// 用户输入的没有加密的密码
			password:   "123456*hello",
			wantedUser: domain.User{},
			wantedErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, errors.New("DB错误"))
				return repo
			},
			//ctx:   context.Background(),
			email: "123@qq.com",
			// 用户输入的没有加密的密码
			password:   "123456*hello",
			wantedUser: domain.User{},
			wantedErr:  errors.New("DB错误"),
		},
		{
			name: "密码不对",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email: "123@qq.com",
					// 应该是一个加密后正确的密码
					Password: "$2a$10$FgqaFMtlGQomZU9Pi/26MeFGTVOrg0bjDzjcu2c/F6ZBx.MQoEP/O",
					Phone:    "12345678",
				}, nil)
				return repo
			},
			//ctx:   context.Background(),
			email: "123@qq.com",
			// 用户输入的没有加密的密码
			password:   "123456*helldkjk",
			wantedUser: domain.User{},
			wantedErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := tc.mock(ctrl)
			userSvc := NewUserService(userRepo)

			user, err := userSvc.Login(tc.ctx, tc.email, tc.password)

			assert.Equal(t, tc.wantedErr, err)
			assert.Equal(t, tc.wantedUser, user)
		})
	}
}
