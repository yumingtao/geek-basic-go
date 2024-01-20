package service

import (
	"context"
	"geek-basic-go/webook/internal/repository"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
}

type InteractiveServiceImpl struct {
	repo repository.InteractiveRepository
}

func NewInteractiveServiceImpl(repo repository.InteractiveRepository) InteractiveService {
	return &InteractiveServiceImpl{repo: repo}
}

func (i *InteractiveServiceImpl) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return i.repo.IncrReadCnt(ctx, biz, bizId)
}
