package repository

import (
	"context"
	"errors"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository/cache"
	"geek-basic-go/webook/internal/repository/dao"
	"geek-basic-go/webook/pkg/logger"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
}

type CachedInteractiveRepository struct {
	dao   dao.InteractiveDao
	cache cache.InteractiveCache
	l     logger.LoggerV1
}

func NewCachedInteractiveRepository(dao dao.InteractiveDao,
	l logger.LoggerV1, cache cache.InteractiveCache) InteractiveRepository {
	return &CachedInteractiveRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

func (c *CachedInteractiveRepository) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	intr, err := c.cache.Get(ctx, biz, id)
	if err == nil {
		return intr, nil
	}
	ie, err := c.dao.Get(ctx, biz, id)
	if err != nil {
		return domain.Interactive{}, err
	}
	intr = c.toDomain(ie)
	// 回写缓存
	err = c.cache.Set(ctx, biz, id, intr)
	if err != nil {
		c.l.Error("回写缓存失败",
			logger.String("biz", biz),
			logger.Int64("bizId", id),
			logger.Error(err),
		)
	}
	return intr, err
}

func (c *CachedInteractiveRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, biz, id, uid)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, dao.ErrRecordNotFound):
		return false, nil
	default:
		return false, err
	}
}

func (c *CachedInteractiveRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetCollectInfo(ctx, biz, id, uid)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, dao.ErrRecordNotFound):
		return false, nil
	default:
		return false, err
	}
}

func (c *CachedInteractiveRepository) AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	err := c.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		Uid:   uid,
		Biz:   biz,
		BizId: id,
		Cid:   cid,
	})
	if err != nil {
		return err
	}
	return c.cache.IncrCollectCntIfPresent(ctx, biz, id)
}

func (c *CachedInteractiveRepository) IncrLike(ctx context.Context, biz string, id int64, uid int64) error {
	err := c.dao.InsertLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return c.cache.IncrLikeCntIfPresent(ctx, biz, id)
}

func (c *CachedInteractiveRepository) DecrLike(ctx context.Context, biz string, id int64, uid int64) error {
	err := c.dao.DeleteLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return c.cache.DecrLikeCntIfPresent(ctx, biz, id)
}

func (c *CachedInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := c.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	// 更新缓存
	// 如果更新失败，造成数据不一致，但影响不大
	return c.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}

func (c *CachedInteractiveRepository) toDomain(ie dao.Interactive) domain.Interactive {
	return domain.Interactive{
		ReadCnt:    ie.ReadCnt,
		LikeCnt:    ie.LikeCnt,
		CollectCnt: ie.CollectCnt,
	}
}
