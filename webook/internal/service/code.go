package service

import (
	"context"
	"errors"
	"fmt"
	"geek-basic-go/webook/internal/repository"
	"geek-basic-go/webook/internal/service/sms"
	"math/rand"
)

var ErrCodeSentTooMany = repository.ErrorCodeSentTooMany

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CodeServiceImpl struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &CodeServiceImpl{
		repo: repo,
		sms:  smsSvc,
	}
}

func (svc *CodeServiceImpl) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, biz, phone, code)
	// 开始发送验证码
	if err != nil {
		return err
	}
	const codeTplId = "1877556"
	return svc.sms.Send(ctx, codeTplId, []string{code}, phone)
}

func (svc *CodeServiceImpl) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {

	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if errors.Is(err, repository.ErrCodeVerifyTooMany) {
		// 对外屏蔽了验证次数过多的错误，直接告诉这个不对
		return false, nil
	}
	return ok, nil
}

func (svc *CodeServiceImpl) generate() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
