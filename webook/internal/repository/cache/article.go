package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"geek-basic-go/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error
	DeleteFirstPage(ctx context.Context, uid int64) error
}

type ArticleRedisCache struct {
	client redis.Cmdable
}

func (a *ArticleRedisCache) DeleteFirstPage(ctx context.Context, uid int64) error {
	key := a.firstKey(uid)
	return a.client.Del(ctx, key).Err()
}

func (a *ArticleRedisCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	key := a.firstKey(uid)
	//val, err := a.client.Get(ctx, firstKey).Result()
	val, err := a.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var res []domain.Article
	//err = json.Unmarshal([]byte(val), &res)
	err = json.Unmarshal(val, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *ArticleRedisCache) SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error {
	// Content只缓存摘要
	for i := range arts {
		arts[i].Content = arts[i].Abstract()
	}
	key := a.firstKey(uid)
	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return a.client.Set(ctx, key, val, time.Minute*10).Err()
}

func (a *ArticleRedisCache) firstKey(uid int64) string {
	return fmt.Sprintf("article:first_page:%d", uid)
}
