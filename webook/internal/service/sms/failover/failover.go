package failover

import (
	"context"
	"errors"
	"geek-basic-go/webook/internal/service/sms"
	"log"
	"sync/atomic"
)

type SupportFailoverSmsService struct {
	svcs []sms.Service
	// v1的字段
	// 当前服务商下标
	idx uint64
}

func NewSupportFailoverSmsService(svcs []sms.Service) *SupportFailoverSmsService {
	return &SupportFailoverSmsService{
		svcs: svcs,
	}
}

func (s *SupportFailoverSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for _, svc := range s.svcs {
		err := svc.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return nil
		}
		log.Println(err)
	}
	return errors.New("发送失败，所有的服务上都尝试过了")
}

// SendV1 is a method of SupportFailoverSmsService that sends SMS messages using a failover mechanism.
// It iterates through a list of SMS services and tries to send the message using each service until a successful response is received.
// If a service returns an error indicating that the context was canceled or the deadline was exceeded, it moves to the next service.
// If all services have been tried and none were successful, it returns an error indicating that the message sending failed on all services.
// 从起始下标开始轮询
// 并且出错也轮询
func (s *SupportFailoverSmsService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&s.idx, 1)
	length := uint64(len(s.svcs))
	for i := idx; i < idx+length; i++ {
		svc := s.svcs[i%length]
		err := svc.Send(ctx, tplId, args, numbers...)
		switch {
		case err == nil:
			return nil
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			// 前者是被调用者取消
			// 后者是调用者设置的超时时间到了
			return err
		}
		log.Println(err)
	}
	return errors.New("发送失败，所有的服务上都尝试过了")
}
