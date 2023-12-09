package service

import (
	"context"
	"errors"
	"fmt"
	"geek-basic-go/webook/internal/repository"
	"geek-basic-go/webook/internal/service/sms"
	"math/rand"
)

type CodeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, biz, phone, code)
	// 开始发送验证码
	if err != nil {
		return err
	}
	const codeTplId = "1877556"
	return svc.sms.Send(ctx, codeTplId, []string{code}, phone)
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {

	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if errors.Is(err, repository.ErrCodeVerifyTooMany) {
		// 对外屏蔽了验证次数过多的错误，直接告诉这个不对
		return false, nil
	}
	return ok, nil
}

func (svc *CodeService) generate() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
