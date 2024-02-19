package cache

import (
	"context"
	"encoding/json"
	"geek-basic-go/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type TopLikedCache interface {
	Set(ctx context.Context, arts []domain.Article) error
	Get(ctx context.Context) ([]domain.Article, error)
}

type TopLikedRedisCache struct {
	client     redis.Cmdable
	key        string
	expiration time.Duration
}

func NewTopLikedRedisCache(client redis.Cmdable) TopLikedCache {
	return &TopLikedRedisCache{
		client:     client,
		key:        "article:top_n:like",
		expiration: time.Minute * 15,
	}
}

func (t *TopLikedRedisCache) Set(ctx context.Context, arts []domain.Article) error {
	for i := range arts {
		arts[i].Content = arts[i].Abstract()
	}
	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return t.client.Set(ctx, t.key, val, t.expiration).Err()
}

func (t *TopLikedRedisCache) Get(ctx context.Context) ([]domain.Article, error) {
	val, err := t.client.Get(ctx, t.key).Bytes()
	if err != nil {
		return nil, err
	}
	var arts []domain.Article
	err = json.Unmarshal(val, &arts)
	return arts, nil
}
