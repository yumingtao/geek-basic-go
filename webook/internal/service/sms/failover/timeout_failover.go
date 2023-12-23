package failover

import (
	"context"
	"errors"
	"geek-basic-go/webook/internal/service/sms"
	"sync/atomic"
)

// 连续N个超时就切换

type TimeoutFailoverSmsService struct {
	svcs []sms.Service
	// 当前正在使用的节点
	idx int32
	// 连续超时计数
	cnt int32
	// 切换阈值，只读，没有被任何地发改的数据是天然线程安全的
	threshold int32
}

func NewTimeoutFailoverSmsService(svcs []sms.Service, threshold int32) *TimeoutFailoverSmsService {
	return &TimeoutFailoverSmsService{
		svcs:      svcs,
		threshold: threshold,
	}
}

func (t *TimeoutFailoverSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)

	// 超过阈值，进行切换
	if cnt >= t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// 切换成功，重置计数
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = newIdx
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tplId, args, numbers...)
	switch {
	case err == nil:
		// err为nil，表明不超时了，需要将cnt重置为0
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	case errors.Is(err, context.DeadlineExceeded):
		// 后者是调用者设置的超时时间到了
		// 超时了，增加计数
		atomic.AddInt32(&t.cnt, 1)
	default:
		// 其他错误，需要考虑怎么处理
		// 可以增加计数，也可以不增加
		// 如果强调超时那么可以不增加
		//atomic.AddInt32(&t.cnt, 1)
		// 如果是EOF之类的错误，直接考虑切换
	}
	return err
}
