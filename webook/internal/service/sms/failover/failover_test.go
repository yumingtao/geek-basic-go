package failover

import (
	"context"
	"errors"
	"geek-basic-go/webook/internal/service/sms"
	smsmocks "geek-basic-go/webook/internal/service/sms/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestSupportFailoverSmsService_Send(t *testing.T) {
	testCases := []struct {
		name      string
		mocks     func(ctrl *gomock.Controller) []sms.Service
		wantedErr error
	}{
		{
			name: "一次发送成功",
			mocks: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0}
			},
			wantedErr: nil,
		},
		{
			name: "第二次发送成功",
			mocks: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("第一次发送失败"))
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return []sms.Service{svc0, svc1}
			},
			wantedErr: nil,
		},
		{
			name: "两次都失败",
			mocks: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("第一次发送失败"))
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("第二次发送失败"))
				return []sms.Service{svc0, svc1}
			},
			wantedErr: errors.New("发送失败，所有的服务上都尝试过了"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			scv := NewSupportFailoverSmsService(tc.mocks(ctrl))
			err := scv.Send(context.Background(), "123", []string{"456"}, "123456789")
			assert.Equal(t, tc.wantedErr, err)
		})
	}
}
