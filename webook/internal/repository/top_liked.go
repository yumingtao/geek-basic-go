package repository

import (
	"context"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository/cache"
	"geek-basic-go/webook/internal/repository/dao"
	"geek-basic-go/webook/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
	rlock "github.com/gotomicro/redis-lock"
	"time"
)

type TopLikedRepository interface {
	GetTopN(ctx context.Context, n int) ([]domain.Article, error)
}

type CachedTopLikedRepository struct {
	intrDao    dao.InteractiveDao
	artDao     dao.ArticleDao
	redisCache cache.TopLikedCache
	localCache *cache.TopLikedLocalCache
	client     *rlock.Client
	key        string
	expiration time.Duration
	l          logger.LoggerV1
}

func NewCachedTopLikedRepository(intrDao dao.InteractiveDao, artDao dao.ArticleDao,
	redisCache cache.TopLikedCache, localCache *cache.TopLikedLocalCache,
	client *rlock.Client, l logger.LoggerV1) TopLikedRepository {
	return &CachedTopLikedRepository{
		intrDao:    intrDao,
		artDao:     artDao,
		redisCache: redisCache,
		localCache: localCache,
		client:     client,
		key:        "lock:top_n:like",
		expiration: time.Second * 30,
		l:          l,
	}
}

func (c *CachedTopLikedRepository) GetTopN(ctx context.Context, n int) ([]domain.Article, error) {
	res, err := c.localCache.Get(ctx)
	// 如果没有err，直接返回本地缓存数据
	if err == nil {
		return res, nil
	}
	// 本地缓存失效，获取Redis缓存
	res, err = c.redisCache.Get(ctx)
	if err == nil {
		return res, nil
	}

	// Redis缓存失效，获取分布式锁并进行数据库计算
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*40)
	defer cancel()
	lock, err := c.client.TryLock(ctx, c.key, c.expiration)
	if err != nil {
		//获取分布式锁失败，直接返回本地缓存数据
		return c.localCache.ForceGet(ctx)
	}

	// 获取分布式锁成功，进行数据库计算
	intrs, err := c.intrDao.GetTopLikedList(ctx, n)
	if err != nil {
		//查询数据库失败，直接返回本地缓存数据
		return c.localCache.ForceGet(ctx)
	}
	ids := make([]int64, 0, len(intrs))
	for _, intr := range intrs {
		ids = append(ids, intr.BizId)
	}
	arts, err := c.artDao.GetPubByIds(ctx, ids)
	if err != nil {
		//查询数据库失败，直接返回本地缓存数据
		return c.localCache.ForceGet(ctx)
	}
	res = slice.Map[dao.Article, domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		return c.toDomain(src)
	})
	// 回写缓存，先更新本地缓存，再更新Redis缓存
	err = c.localCache.Set(ctx, res)
	if err != nil {
		c.l.Error("存入本地缓存失败", logger.Error(err))
	}
	err = c.redisCache.Set(ctx, res)
	if err != nil {
		c.l.Error("存入Redis缓存失败", logger.Error(err))
	}

	// 释放分布式锁
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err = lock.Unlock(ctx)
	if err != nil {
		// 释放分布式锁失败不影响正常业务，锁过期了也会释放，所以这里只需要记录日志，不需要返回nil和err
		c.l.Error("释放Redis分布式锁失败", logger.Error(err))
		//return nil, err
	}
	return res, nil
}

func (c *CachedTopLikedRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Status: domain.ArticleStatus(art.Status),
		Ctime:  time.UnixMilli(art.Ctime),
		Utime:  time.UnixMilli(art.Utime),
	}
}
