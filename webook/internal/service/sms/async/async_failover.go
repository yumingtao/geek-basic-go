package async

import (
	"context"
	"errors"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository"
	"geek-basic-go/webook/internal/service/sms"
	"log"
	"time"
)

type AsyncFailoverSmsService struct {
	svc sms.Service
	reo repository.SmsRepository
}

func (a *AsyncFailoverSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	isNeedAsyncSend := a.isNeedAsyncSend()
	if isNeedAsyncSend {
		// 将该条sms加入异步发送数据库
		s := domain.Sms{
			Tpl:         tplId,
			Args:        args,
			Numbers:     numbers,
			RetryMaxCnt: 3,
		}
		err := a.reo.Add(ctx, s)
		if err != nil {
			return err
		}
	}
	// 不需要异步发送
	err := a.svc.Send(ctx, tplId, args, numbers...)
	return err
}

func (a *AsyncFailoverSmsService) isNeedAsyncSend() bool {
	//实现判断是否需要异步发送的逻辑
	// 1. 响应时间 2.
	return true
}

// StartAsync 启动一个无限循环，每次循环执行异步发送操作
func (a *AsyncFailoverSmsService) StartAsync() {
	time.Sleep(time.Second * 3)
	for {
		a.AsyncSend()
	}
}

func (a *AsyncFailoverSmsService) AsyncSend() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	s, err := a.reo.GetEarliestSmsByInterval(ctx, time.Minute)
	cancel()
	switch {
	case err == nil:
		//找到了需要异步发送的sms
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := a.svc.Send(ctx, s.Tpl, s.Args, s.Numbers...)
		if err != nil {
			log.Printf("异步发送sms:%d失败", s.Id)
		}
		status := err == nil
		// 更新本次异步发送结果
		err = a.reo.MarkStatus(ctx, s.Id, status)
		if err != nil {
			log.Printf("异步发送sms:%d成功, 标记数据库状态失败", s.Id)
		}
	case errors.Is(err, repository.ErrSmsNotFound):
		// 没有需要执行异步发送的sms，sleep一段时间之后再查询
		time.Sleep(time.Second * 3)
	default:
		// 执行数据库查询失败，sleep一段时间之后再查询
		log.Println("查询异步短信发送数据库失败")
		time.Sleep(time.Second * 3)
	}
}
