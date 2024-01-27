package service

import (
	"context"
	"geek-basic-go/webook/internal/repository"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, id int64, uid int64) error
	CancelLike(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
}

type InteractiveServiceImpl struct {
	repo repository.InteractiveRepository
}

func NewInteractiveServiceImpl(repo repository.InteractiveRepository) InteractiveService {
	return &InteractiveServiceImpl{repo: repo}
}

func (i *InteractiveServiceImpl) Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	return i.repo.AddCollectionItem(ctx, biz, id, cid, uid)
}

func (i *InteractiveServiceImpl) Like(ctx context.Context, biz string, id int64, uid int64) error {
	return i.repo.IncrLike(ctx, biz, id, uid)
}

func (i *InteractiveServiceImpl) CancelLike(ctx context.Context, biz string, id int64, uid int64) error {
	return i.repo.DecrLike(ctx, biz, id, uid)
}

func (i *InteractiveServiceImpl) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return i.repo.IncrReadCnt(ctx, biz, bizId)
}
