package cache

import (
	"context"
	"errors"
	"geek-basic-go/webook/internal/domain"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"time"
)

type TopLikedLocalCache struct {
	topN       *atomicx.Value[[]domain.Article]
	ddl        *atomicx.Value[time.Time]
	expiration time.Duration
}

func NewTopLikedLocalCache() *TopLikedLocalCache {
	return &TopLikedLocalCache{
		topN:       atomicx.NewValue[[]domain.Article](),
		ddl:        atomicx.NewValue[time.Time](),
		expiration: time.Minute * 30,
	}
}

func (t *TopLikedLocalCache) Set(ctx context.Context, arts []domain.Article) error {
	t.topN.Store(arts)
	t.ddl.Store(time.Now().Add(t.expiration))
	return nil
}

func (t *TopLikedLocalCache) Get(ctx context.Context) ([]domain.Article, error) {
	ddl := t.ddl.Load()
	arts := t.topN.Load()
	if len(arts) == 0 || ddl.Before(time.Now()) {
		return nil, errors.New("本地缓存失效")
	}
	return arts, nil
}

func (t *TopLikedLocalCache) ForceGet(ctx context.Context) ([]domain.Article, error) {
	arts := t.topN.Load()
	if len(arts) == 0 {
		return nil, errors.New("本地缓存失效")
	}
	return arts, nil
}
