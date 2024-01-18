package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository/dao"
	"github.com/redis/go-redis/v9"
	"time"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error
	DeleteFirstPage(ctx context.Context, uid int64) error
	Get(ctx context.Context, id int64) (domain.Article, error)
	Set(ctx context.Context, art dao.Article) error
	GetPub(ctx context.Context, id int64) (domain.Article, error)
	SetPub(ctx context.Context, res domain.Article) error
}

type ArticleRedisCache struct {
	client redis.Cmdable
}

func NewArticleRedisCache(client redis.Cmdable) ArticleCache {
	return &ArticleRedisCache{
		client: client,
	}
}

func (a *ArticleRedisCache) SetPub(ctx context.Context, art domain.Article) error {
	val, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return a.client.Set(ctx, a.pubKey(art.Id), val, time.Minute*10).Err()
}

func (a *ArticleRedisCache) GetPub(ctx context.Context, id int64) (domain.Article, error) {
	val, err := a.client.Get(ctx, a.pubKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var res domain.Article
	err = json.Unmarshal(val, &res)
	return res, err
}

func (a *ArticleRedisCache) Get(ctx context.Context, id int64) (domain.Article, error) {
	val, err := a.client.Get(ctx, a.key(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var res domain.Article
	err = json.Unmarshal(val, &res)
	return res, err
}

func (a *ArticleRedisCache) Set(ctx context.Context, art dao.Article) error {
	// JSON 序列化大部分场景，都不会引起性能问题
	// 如果有性能问题,可以换一个序列化方法，可以使用protobuf，也可以考虑别的，例如 gob(Go Object Binary)
	val, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return a.client.Set(ctx, a.key(art.Id), val, time.Minute*1).Err()
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

func (a *ArticleRedisCache) key(id int64) string {
	return fmt.Sprintf("article:detail:%d", id)
}

func (a *ArticleRedisCache) pubKey(id int64) string {
	return fmt.Sprintf("article:pub:detail:%d", id)
}
