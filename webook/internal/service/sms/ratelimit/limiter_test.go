package ratelimit

import (
	"context"
	"errors"
	"geek-basic-go/webook/internal/service/sms"
	smsmocks "geek-basic-go/webook/internal/service/sms/mocks"
	"geek-basic-go/webook/pkg/limiter"
	limitermocks "geek-basic-go/webook/pkg/limiter/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestLimitSmsService_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter)
		// 一个输入都没有
		wantedErr error
	}{
		{
			name: "不限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return svc, l
			},
			wantedErr: nil,
		},
		{
			name: "限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				//svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return svc, l
			},
			wantedErr: errLimited,
		},
		{
			name: "限流器错误",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, errors.New("redis限流器错误"))
				//svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return svc, l
			},
			wantedErr: errors.New("redis限流器错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			smsSvc, l := tc.mock(ctrl)
			svc := NewLimitSmsService(smsSvc, l)
			err := svc.Send(context.Background(), "abc", []string{"123"}, "13871173212")
			assert.Equal(t, tc.wantedErr, err)
		})
	}
}
