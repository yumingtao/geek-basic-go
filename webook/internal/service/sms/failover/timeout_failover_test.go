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

func Test_timeoutFailoverSmsService_Send(t *testing.T) {
	testCases := []struct {
		name      string
		mocks     func(ctrl *gomock.Controller) []sms.Service
		idx       int32
		cnt       int32
		threshold int32
		wantedErr error
		wantedIdx int32
		wantedCnt int32
	}{
		{
			name: "没有切换",
			mocks: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0}
			},
			idx:       0,
			cnt:       2,
			threshold: 3,
			wantedErr: nil,
			wantedIdx: 0,
			// 成功后，重置了超时计数
			wantedCnt: 0,
		},
		{
			name: "触发切换，成功",
			mocks: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0, svc1}
			},
			idx:       0,
			cnt:       3,
			threshold: 3,
			wantedErr: nil,
			wantedIdx: 1,
			wantedCnt: 0,
		},
		{
			name: "触发切换，失败",
			mocks: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("切换失败"))
				return []sms.Service{svc0, svc1}
			},
			// 原来是1，出发了切换，然后切换到0，但是svc0在发送的时候失败了
			idx:       1,
			cnt:       3,
			threshold: 3,
			wantedErr: errors.New("切换失败"),
			// 这里应该是0，因为触发了切换
			wantedIdx: 0,
			// 这里也是0，因为触发了切换，但是走到了default分支，default分支什么都没干
			wantedCnt: 0,
		},
		{
			name: "触发切换，超时",
			mocks: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(context.DeadlineExceeded)
				return []sms.Service{svc0, svc1}
			},
			// 原来是1，出发了切换，然后切换到0，但是svc0在发送时超时了
			idx:       1,
			cnt:       3,
			threshold: 3,
			wantedErr: context.DeadlineExceeded,
			// 这里应该是0，因为触发了切换
			wantedIdx: 0,
			// 这里也是1，因为触发了切换后，svc0发送超时，所以自增
			wantedCnt: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewTimeoutFailoverSmsService(tc.mocks(ctrl), tc.threshold)
			svc.idx = tc.idx
			svc.cnt = tc.cnt
			err := svc.Send(context.Background(), "123", []string{"456"}, "123456789")
			assert.Equal(t, tc.wantedErr, err)
			assert.Equal(t, tc.wantedIdx, svc.idx)
			assert.Equal(t, tc.wantedCnt, svc.cnt)
		})
	}

}
